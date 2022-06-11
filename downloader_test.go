package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveTemp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping this SaveTemp is too expensive")
	}

	contents := []byte("Hello, World!")
	path, err := SaveTemp(&contents)

	if err != nil {
		t.Fatal()
	}

	defer func() {
		_ = Afero.Remove(path)
	}()

	if !strings.HasPrefix(filepath.Base(path), TempPrefix) {
		t.Error(TempPrefix + " is expected as a temp file prefix")
	}

	if exists, _ := Afero.Exists(path); !exists {
		t.Fatal("File was not created")
	}

	newContents, err := Afero.ReadFile(path)
	if err != nil {
		t.Fatal("Can not open the file")
	}

	if string(newContents) != string(contents) {
		t.Error("Written contents differ from the original")
	}
}

func TestRemoveIfExists(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping this RemoveIfExists is too expensive")
	}

	contents := []byte("Hello, World!")

	path, _ := SaveTemp(&contents)

	if err := RemoveIfExists(path); err != nil {
		t.Fatal("Error while removing file")
	}

	if exists, _ := Afero.Exists(path); exists {
		t.Error("File was not removed")
	}

	if err := RemoveIfExists(path); err != nil {
		t.Error("Unexpected error when removing file that doesn't exist")
	}

}

func TestDownloadChapter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping this DownloadChapter is too expensive")
	}

	manga, _ := DefaultConfig().Scrapers[0].SearchManga(TestQuery)
	chapters, _ := manga[0].Scraper.GetChapters(manga[0])

	path, err := DownloadChapter(chapters[0], nil, true)

	if err != nil {
		t.Fatal("Chapter was not downloaded")
	}

	defer func() {
		_ = Afero.RemoveAll(filepath.Dir(path))
	}()

	if !strings.HasPrefix(path, os.TempDir()) {
		t.Error("Chapter is expected to be downloaded in the temp folder")
	}

	if !strings.HasPrefix(filepath.Base(path), TempPrefix) {
		t.Error(TempPrefix + " is expected as a temp chapter prefix")
	}

	if exists, _ := Afero.Exists(path); !exists {
		t.Fatal("Chapter was not downloaded")
	}

	if stat, _ := Afero.Stat(path); stat.Size() <= 0 {
		t.Error("Downloaded pdf is empty")
	}
}
