package updater

import (
	"encoding/json"
	"errors"
	"github.com/metafates/mangal/cache"
	"github.com/metafates/mangal/util"
	"net/http"
	"time"
)

var versionCacher = cache.New[string]("version", &cache.Options[string]{
	ExpireEvery: time.Hour,
})

// LatestVersion returns the latest version of mangal.
// It will fetch the latest version from the GitHub API.
func LatestVersion() (version string, err error) {
	ver := versionCacher.Get()
	if ver.IsPresent() {
		return ver.MustGet(), nil
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
