package manager

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/go-git/go-git/v5"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/luaprovider"
	"github.com/mangalorg/mangal/afs"
	"github.com/mangalorg/mangal/path"
	"github.com/mangalorg/mangal/provider/info"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/afero"
)

type (
	AddOptions struct {
		URL *url.URL
	}

	UpdateOptions struct {
	}
)

func Add(ctx context.Context, options AddOptions) error {
	tempDir, err := afs.Afero.TempDir(path.TempDir(), "")
	if err != nil {
		return err
	}

	fmt.Println(tempDir)
	_, err = git.PlainCloneContext(ctx, tempDir, false, &git.CloneOptions{
		URL:      options.URL.String(),
		Progress: os.Stdout, // TODO: change this
	})

	if err != nil {
		return err
	}

	infoFile, err := afs.Afero.OpenFile(filepath.Join(tempDir, info.Filename), os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer infoFile.Close()

	providerInfo, err := info.New(infoFile)
	if err != nil {
		return err
	}

	ID := providerInfo.ID
	if ID == "" {
		return fmt.Errorf("ID is empty")
	}

	loaders, err := Loaders()
	if err != nil {
		return err
	}

	_, IDExists := lo.Find(loaders, func(loader libmangal.ProviderLoader) bool {
		return loader.Info().ID == ID
	})

	if IDExists {
		return fmt.Errorf("provider with ID %q already exists", ID)
	}

	target := filepath.Join(path.ProvidersDir(), ID)
	fmt.Println(target)
	return afs.Afero.Rename(tempDir, target)
}

func Update(ctx context.Context, options UpdateOptions) error {
	providersDir := path.ProvidersDir()
	dirEntries, err := afs.Afero.ReadDir(providersDir)
	if err != nil {
		return err
	}

	for _, dirEntry := range dirEntries {
		repo, err := git.PlainOpen(filepath.Join(providersDir, dirEntry.Name()))

		if errors.Is(err, git.ErrRepositoryNotExists) {
			continue
		}

		if err != nil {
			return err
		}

		worktree, err := repo.Worktree()
		if err != nil {
			return err
		}

		err = worktree.PullContext(ctx, &git.PullOptions{
			Progress: os.Stdout,
			Force:    true,
		})

		if !(errors.Is(err, git.NoErrAlreadyUpToDate) || errors.Is(err, git.ErrRemoteNotFound)) {
			return err
		}
	}

	return nil
}

func Remove(tag string) error {
	return afs.Afero.RemoveAll(filepath.Join(path.ProvidersDir(), tag))
}

type NewOptions struct {
	info.Info
	Dir string
}

func New(options NewOptions) error {
	if !options.Type.IsAType() {
		return fmt.Errorf("invalid provider type: %v", options.Type)
	}

	dir := options.Dir
	providerPath := filepath.Join(dir, options.Info.ID)

	files := afero.Afero{Fs: afero.NewMemMapFs()}
	switch options.Type {
	case info.TypeLua:
		if err := newLua(files, options.Info); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported provider type %s", options.Type)
	}

	err := files.WriteFile("README.md", []byte(options.Markdown()), 0755)
	if err != nil {
		return err
	}

	return createRepo(providerPath, files)
}

func newLua(af afero.Afero, information info.Info) error {
	err := af.WriteFile(".gitignore", []byte("sdk.lua"), 0755)
	if err != nil {
		return err
	}

	infoFile, err := af.Create(info.Filename)
	if err != nil {
		return err
	}
	defer infoFile.Close()

	err = toml.NewEncoder(infoFile).Encode(information)
	if err != nil {
		return err
	}

	err = af.WriteFile("main.lua", []byte(luaprovider.LuaTemplate()), 0755)
	if err != nil {
		return err
	}

	return af.WriteFile("sdk.lua", []byte(luaprovider.LuaDoc()), 0755)
}

func createRepo(dir string, files afero.Fs) error {
	err := afero.Walk(files, ".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == "" {
			return nil
		}

		dstPath := filepath.Join(dir, path)
		if info.IsDir() {
			return afs.Afero.MkdirAll(dstPath, info.Mode())
		}
		src, err := files.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := afs.Afero.Create(dstPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		return err
	})

	if err != nil {
		return err
	}

	repo, err := git.PlainInit(dir, false)
	if err != nil {
		return err
	}

	tree, err := repo.Worktree()
	if err != nil {
		return err
	}

	return tree.AddWithOptions(&git.AddOptions{
		All: true,
	})
}
