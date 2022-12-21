package version

import (
	"encoding/json"
	"errors"
	"github.com/metafates/gache"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"net/http"
	"path/filepath"
	"time"
)

var versionCacher = gache.New[string](&gache.Options{
	Path:       filepath.Join(where.Cache(), "version.json"),
	Lifetime:   time.Hour * 24 * 2,
	FileSystem: &filesystem.GacheFs{},
})

// Latest returns the latest version of mangal.
// It will fetch the latest version from the GitHub API.
func Latest() (version string, err error) {
	ver, expired, err := versionCacher.Get()
	if err != nil {
		return "", err
	}

	if !expired && ver != "" {
		return ver, nil
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
	_ = versionCacher.Set(version)
	return
}
