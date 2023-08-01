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
)

func Add(ctx context.Context, tag string, URL *url.URL) error {
	providerPath := filepath.Join(path.ProvidersDir(), tag)

	exists, err := fs.Afero.Exists(providerPath)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("provider with tag %q already exists", tag)
	}

	_, err = git.PlainCloneContext(ctx, providerPath, false, &git.CloneOptions{
		URL:      URL.String(),
		Progress: os.Stdout, // TODO: change this
	})

	if err != nil {
		return err
	}

	return nil

	//worktree, err := repo.Worktree()
	//if err != nil {
	//	return err
	//}
	//
	//pullOptions := &git.PullOptions{
	//	Progress: os.Stdout,
	//	Force:    true,
	//}
	//
	//err = worktree.PullContext(ctx, pullOptions)
	//
	//if !errors.Is(err, git.NoErrAlreadyUpToDate) {
	//	return err
	//}
	//
	//if provider.Rev == "" {
	//	return nil
	//}
	//
	//return worktree.Checkout(&git.CheckoutOptions{
	//	Hash:  plumbing.Hash([]byte(provider.Rev)),
	//	Force: true,
	//})
}
