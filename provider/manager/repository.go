package manager

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"github.com/pkg/errors"
)

type (
	AddOptions struct {
		Tag string
		URL *url.URL
	}

	UpdateOptions struct {
		Tag string
	}
)

func Add(ctx context.Context, options AddOptions) error {
	providerPath := filepath.Join(path.ProvidersDir(), options.Tag)

	exists, err := fs.Afero.Exists(providerPath)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("provider with tag %q already exists", options.Tag)
	}

	_, err = git.PlainCloneContext(ctx, providerPath, false, &git.CloneOptions{
		URL:      options.URL.String(),
		Progress: os.Stdout, // TODO: change this
	})

	if err != nil {
		return err
	}

	return nil
}

func Update(ctx context.Context, options UpdateOptions) error {
	providersDir := path.ProvidersDir()
	dirEntries, err := fs.Afero.ReadDir(providersDir)
	if err != nil {
		return err
	}

	for _, dirEntry := range dirEntries {
		if options.Tag != "" && dirEntry.Name() != options.Tag {
			continue
		}

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

		if options.Tag != "" {
			break
		}
	}

	return nil
}

func Remove(tag string) error {
	return fs.Afero.RemoveAll(filepath.Join(path.ProvidersDir(), tag))
}
