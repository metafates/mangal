package scraper

import (
	"encoding/json"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/util"
	"github.com/spf13/afero"
	"testing"
)

func TestAnilistClient_SearchManga(t *testing.T) {
	// sample anilist client
	anilistClient, err := NewAnilistClient("", "")

	if err != nil {
		t.Error(err)
	}

	// search for manga
	manga, err := anilistClient.SearchManga(common.TestQuery)

	// check for error
	if err != nil {
		t.Error(err)
	}

	// check for results
	if len(manga) == 0 {
		t.Error("no results found")
	}
}

func TestAnilistClient_SavePreferences(t *testing.T) {
	// sample anilist client
	anilistClient, err := NewAnilistClient("", "")

	if err != nil {
		t.Error(err)
	}

	// save preferences
	err = anilistClient.SavePreferences()

	// check for error
	if err != nil {
		t.Error(err)
	}

	// check if file was created
	anilistFile, err := util.AnilistCacheFile()

	if err != nil {
		t.Error(err)
	}

	if exists, err := afero.Exists(filesystem.Get(), anilistFile); err != nil {
		t.Error(err)
	} else if !exists {
		t.Error("file was not created")
	}

	// check if file contains correct data
	file, err := afero.ReadFile(filesystem.Get(), anilistFile)
	if err != nil {
		t.Error(err)
	}

	var preferences AnilistPreferences
	err = json.Unmarshal(file, &preferences)
	if err != nil {
		t.Error(err)
	}

	// compare preferences
	if preferences.Token != anilistClient.Preferences.Token {
		t.Error("preferences were not saved")
	}

	if len(preferences.Connections) != len(anilistClient.Preferences.Connections) {
		t.Error("preferences were not saved")
	}

	// delete file
	if err := filesystem.Get().Remove(anilistFile); err != nil {
		t.Error(err)
	}
}

func TestAnilistClient_ToAnilistURL(t *testing.T) {
	// sample anilist client
	anilistClient, err := NewAnilistClient("", "")

	if err != nil {
		t.Error(err)
	}

	// search manga with default scraper
	conf := config.GetConfig("")
	manga, err := conf.Scrapers[0].SearchManga(common.TestQuery)

	if err != nil {
		return
	}

	// create url
	url := anilistClient.ToAnilistURL(manga[0])

	// check for error
	if err != nil {
		t.Error(err)
	}

	// check for url
	if url == nil {
		t.Error("url was not created")
	}
}
