package main

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Error(err)
	}

	anilistFile := filepath.Join(configDir, Mangal, "anilist.json")

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
