package lua

import (
	"fmt"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/luaprovider"
	"github.com/mangalorg/mangal/cache/bbolt"
	"github.com/mangalorg/mangal/config"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"github.com/philippgille/gokv/encoding"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	infoYAML = "info.yaml"
	mainLua  = "main.lua"
)

func newLoader(dir string) (libmangal.ProviderLoader, error) {
	infoFilePath := filepath.Join(dir, infoYAML)

	exists, err := fs.Afero.Exists(infoFilePath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("%s is missing", infoFilePath)
	}

	infoFile, err := fs.Afero.OpenFile(infoFilePath, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}

	var info libmangal.ProviderInfo

	if err := yaml.NewDecoder(infoFile).Decode(&info); err != nil {
		return nil, err
	}
	infoFile.Close()

	providerMainFilePath := filepath.Join(dir, mainLua)
	exists, err = fs.Afero.Exists(providerMainFilePath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("%s is missing", providerMainFilePath)
	}

	providerMainFileContents, err := fs.Afero.ReadFile(providerMainFilePath)
	if err != nil {
		return nil, err
	}

	ttl, err := time.ParseDuration(config.Config.Cache.Providers.Lua.TTL.Get())
	if err != nil {
		log.Fatal(err)
	}

	store, err := bbolt.NewStore(bbolt.Options{
		TTL:        ttl,
		BucketName: info.Name,
		Path:       filepath.Join(path.CacheDir(), info.Name+".db"),
		Codec:      encoding.Gob,
	})
	if err != nil {
		return nil, err
	}

	options := luaprovider.Options{
		PackagePaths: []string{dir},
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		HTTPStore: store,
	}

	return luaprovider.NewLoader(providerMainFileContents, info, options)
}

func InstalledProviders() ([]libmangal.ProviderLoader, error) {
	dir := path.LuaProvidersDir()
	dirEntries, err := fs.Afero.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var loaders []libmangal.ProviderLoader
	for _, dirEntry := range dirEntries {
		// skip non directories
		if !dirEntry.IsDir() {
			continue
		}

		loader, err := newLoader(filepath.Join(dir, dirEntry.Name()))
		if err != nil {
			return nil, err
		}

		loaders = append(loaders, loader)
	}

	return loaders, nil
}
