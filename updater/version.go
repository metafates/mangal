package updater

import (
	"encoding/json"
	"errors"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	cachedLatestVersion string
	versionCacheFile    = filepath.Join(where.Cache(), "version.json")
)

type versionCache struct {
	Version string    `json:"version"`
	Updated time.Time `json:"updated"`
}

func getCachedVersion() (version string, err error) {
	if cachedLatestVersion != "" {
		return cachedLatestVersion, nil
	}

	exists, err := filesystem.Api().Exists(versionCacheFile)
	if err != nil {
		return
	}

	if !exists {
		return
	}

	var data []byte
	data, err = filesystem.Api().ReadFile(versionCacheFile)
	if err != nil {
		return
	}

	var cache versionCache

	err = json.Unmarshal(data, &cache)
	if err != nil {
		return
	}

	if time.Since(cache.Updated) > time.Hour {
		return
	}

	version = cache.Version
	return
}

func cacheVersion(version string) error {
	cachedLatestVersion = version
	cache := versionCache{
		Version: version,
		Updated: time.Now(),
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return filesystem.Api().WriteFile(versionCacheFile, data, os.ModePerm)
}

// LatestVersion returns the latest version of mangal.
// It will fetch the latest version from the GitHub API.
func LatestVersion() (version string, err error) {
	version, err = getCachedVersion()
	if err == nil && version != "" {
		return
	}

	resp, err := http.Get("https://api.github.com/repos/metafates/mangal/releases/latest")
	if err != nil {
		return
	}

	defer util.Ignore(resp.Body.Close)

	var release struct {
		TagName string `json:"tag_name"`
	}

	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return
	}

	// remove the v from the tag name
	if release.TagName == "" {
		err = errors.New("empty tag name")
		return
	}

	version = release.TagName[1:]
	_ = cacheVersion(version)
	return
}
