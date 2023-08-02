package manager

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"github.com/mangalorg/mangal/provider/info"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type (
	AddOptions struct {
		URL *url.URL
	}

	UpdateOptions struct {
	}
)

func Add(ctx context.Context, options AddOptions) error {
	tempDir, err := fs.Afero.TempDir(path.TempDir(), "")
	if err != nil {
		return err
	}

	_, err = git.PlainCloneContext(ctx, tempDir, false, &git.CloneOptions{
		URL:      options.URL.String(),
		Progress: os.Stdout, // TODO: change this
	})

	if err != nil {
		return err
	}

	infoFilePath := filepath.Join(tempDir, info.Filename)

	infoFile, err := fs.Afero.OpenFile(infoFilePath, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer infoFile.Close()

	providerInfo, err := info.Parse(infoFile)
	if err != nil {
		return err
	}

	ID := providerInfo.Info.ID

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

	return fs.Afero.Rename(tempDir, filepath.Join(path.ProvidersDir(), ID))
}

func Update(ctx context.Context, options UpdateOptions) error {
	providersDir := path.ProvidersDir()
	dirEntries, err := fs.Afero.ReadDir(providersDir)
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
	return fs.Afero.RemoveAll(filepath.Join(path.ProvidersDir(), tag))
}
