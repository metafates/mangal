package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"github.com/mangalorg/mangal/provider/info"
	"github.com/mangalorg/mangal/provider/lua"
	"github.com/pelletier/go-toml"
)

func getLoader(dir string) (libmangal.ProviderLoader, error) {
	infoFile, err := fs.Afero.OpenFile(
		filepath.Join(dir, info.Filename),
		os.O_RDONLY,
		0755,
	)

	if err != nil {
		return nil, err
	}

	defer infoFile.Close()

	decoder := toml.NewDecoder(infoFile)
	decoder.Strict(true)

	var providerInfo info.Info

	if err := decoder.Decode(&providerInfo); err != nil {
		return nil, err
	}

	switch providerInfo.Type {
	case info.TypeLua:
		return lua.NewLoader(providerInfo.Info, dir)
	default:
		return nil, fmt.Errorf("unkown provider type: %#v", providerInfo.Type)
	}
}

func InstalledProviders() ([]libmangal.ProviderLoader, error) {
	providersDir := path.ProvidersDir()
	dirEntries, err := fs.Afero.ReadDir(providersDir)
	if err != nil {
		return nil, err
	}

	var loaders []libmangal.ProviderLoader
	for _, dirEntry := range dirEntries {
		// skip non directories
		if !dirEntry.IsDir() {
			continue
		}

		loader, err := getLoader(filepath.Join(providersDir, dirEntry.Name()))
		if err != nil {
			return nil, err
		}

		loaders = append(loaders, loader)
	}

	return loaders, nil
}
