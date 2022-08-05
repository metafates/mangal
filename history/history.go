package history

import (
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/constants"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"os"
	"path/filepath"
)

type SavedChapter struct {
	SourceID  string `json:"source_id"`
	MangaName string `json:"manga_name"`
	MangaURL  string `json:"manga_url"`
	Name      string `json:"name"`
	URL       string `json:"url"`
}

func Get() (map[string]*SavedChapter, error) {
	historyFile, err := Location()
	if err != nil {
		return nil, err
	}

	// decode json into slice of structs
	var chapters map[string]*SavedChapter
	contents, err := filesystem.Get().ReadFile(historyFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &chapters)
	if err != nil {
		return nil, err
	}

	return chapters, nil
}

func Save(chapter *source.Chapter) error {
	historyFile, err := Location()
	if err != nil {
		return err
	}

	// decode json into slice of structs
	var chapters map[string]*SavedChapter
	contents, err := filesystem.Get().ReadFile(historyFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &chapters)
	if err != nil {
		return err
	}

	jsonChapter := SavedChapter{
		SourceID:  chapter.SourceID,
		MangaName: chapter.Manga.Name,
		MangaURL:  chapter.Manga.URL,
		Name:      chapter.Name,
		URL:       chapter.URL,
	}

	chapters[fmt.Sprintf("%s (%s)", chapter.Manga.Name, chapter.SourceID)] = &jsonChapter

	// encode json
	encoded, err := json.Marshal(chapters)
	if err != nil {
		return err
	}

	// write to file
	err = filesystem.Get().WriteFile(historyFile, encoded, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func Location() (string, error) {
	cacheDir := filepath.Join(lo.Must(os.UserCacheDir()), constants.CachePrefix)
	err := filesystem.Get().MkdirAll(filepath.Dir(cacheDir), os.ModePerm)
	if err != nil {
		return "", err
	}

	historyFile := filepath.Join(cacheDir, constants.History+".json")
	exists, err := filesystem.Get().Exists(historyFile)
	if err != nil {
		return "", err
	}

	if !exists {
		err = filesystem.Get().WriteFile(historyFile, []byte("{}"), os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return historyFile, nil
}
