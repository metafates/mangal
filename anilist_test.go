package main

import (
	"encoding/json"
	"testing"
)

func TestAnilistClient_SearchManga(t *testing.T) {
	// sample anilist client
	anilistClient, err := NewAnilistClient("", "")

	if err != nil {
		t.Error(err)
	}

	// search for manga
	manga, err := anilistClient.SearchManga(TestQuery)

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
	anilistFile, err := AnilistCacheFile()

	if err != nil {
		t.Error(err)
	}

	if exists, err := Afero.Exists(anilistFile); err != nil {
		t.Error(err)
	} else if !exists {
		t.Error("file was not created")
	}

	// check if file contains correct data
	file, err := Afero.ReadFile(anilistFile)
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
	if err := Afero.Remove(anilistFile); err != nil {
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
	conf := GetConfig("")
	manga, err := conf.Scrapers[0].SearchManga(TestQuery)

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
