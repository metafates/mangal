package downloader

import (
	"bytes"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/util"
	"github.com/spf13/afero"
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
	buffer := bytes.NewBuffer(contents)
	path, err := SaveTemp(buffer)

	if err != nil {
		t.Fatal()
	}

	defer func() {
		_ = filesystem.Get().Remove(path)
	}()

	if !strings.HasPrefix(filepath.Base(path), common.TempPrefix) {
		t.Error(common.TempPrefix + " is expected as a temp file prefix")
	}

	if exists, _ := afero.Exists(filesystem.Get(), path); !exists {
		t.Fatal("File was not created")
	}

	newContents, err := afero.ReadFile(filesystem.Get(), path)
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
	buffer := bytes.NewBuffer(contents)

	path, _ := SaveTemp(buffer)

	if err := util.RemoveIfExists(path); err != nil {
		t.Fatal("Error while removing file")
	}

	if exists, _ := afero.Exists(filesystem.Get(), path); exists {
		t.Error("File was not removed")
	}

	if err := util.RemoveIfExists(path); err != nil {
		t.Error("Unexpected error when removing file that doesn't exist")
	}

}

func TestDownloadChapter(t *testing.T) {
	config.Initialize("", false)
	if testing.Short() {
		t.Skip("skipping this DownloadChapter is too expensive")
	}

	manga, _ := config.DefaultConfig().Scrapers[0].SearchManga(common.TestQuery)
	chapters, _ := manga[0].Scraper.GetChapters(manga[0])

	path, err := DownloadChapter(chapters[0], nil, true)

	if err != nil {
		t.Fatal("Chapter was not downloaded")
	}

	defer func() {
		_ = filesystem.Get().RemoveAll(filepath.Dir(path))
	}()

	if !strings.HasPrefix(path, os.TempDir()) {
		t.Error("Chapter is expected to be downloaded in the temp folder")
	}

	if !strings.HasPrefix(filepath.Base(path), common.TempPrefix) {
		t.Error(common.TempPrefix + " is expected as a temp chapter prefix")
	}

	if exists, _ := afero.Exists(filesystem.Get(), path); !exists {
		t.Fatal("Chapter was not downloaded")
	}

	if stat, _ := filesystem.Get().Stat(path); stat.Size() <= 0 {
		t.Error("Downloaded pdf is empty")
	}
}
