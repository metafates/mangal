package updater

import (
	"encoding/json"
	"github.com/metafates/mangal/util"
	"net/http"
)

func LatestVersion() (version string, err error) {
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
	version = release.TagName[1:]
	return
}
