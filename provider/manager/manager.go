package manager

import (
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"github.com/mangalorg/mangal/provider/bundle"
)

func Loaders() ([]libmangal.ProviderLoader, error) {
	return bundle.Loaders(path.ProvidersDir())
}

func Tags() ([]string, error) {
	dirEntries, err := fs.Afero.ReadDir(path.ProvidersDir())
	if err != nil {
		return nil, err
	}

	var tags []string

	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}

		tags = append(tags, dirEntry.Name())
	}

	return tags, nil
}
