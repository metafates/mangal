package history

import (
	"encoding/json"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/scraper"
	"github.com/metafates/mangal/util"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Entry struct {
	Manga      *scraper.URL
	Chapter    *scraper.URL
	SourceName string
}

func (h *Entry) Title() string {
	if config.UserConfig.UI.Icons {
		return "\uF5B9 " + h.Manga.Info
	}
	return h.Manga.Info
}

func (h *Entry) Description() string {
	// replace according to the name description
	description := strings.ReplaceAll(config.UserConfig.UI.ChapterNameTemplate, "%0d", util.PadZeros(h.Chapter.Index, 4))
	description = strings.ReplaceAll(description, "%d", strconv.Itoa(h.Chapter.Index))
	description = strings.ReplaceAll(description, "%s", "\""+h.Chapter.Info+"\"")

	if config.UserConfig.UI.Icons {
		return "\uF129 " + description
	}
	return description
}

func (h *Entry) FilterValue() string {
	return h.Manga.Info
}

func WriteHistory(chapter *scraper.URL) error {
	historyFile, err := util.HistoryFilePath()

	if err != nil {
		return err
	}

	history, err := ReadHistory()

	if err != nil {
		return err
	}

	history[chapter.Relation.Address] = &Entry{
		Manga:      chapter.Relation,
		Chapter:    chapter,
		SourceName: chapter.Scraper.Source.Name,
	}

	historyJSON, err := json.Marshal(history)

	if err != nil {
		return err
	}

	// create the file if it doesn't exist
	if exists, err := afero.Exists(filesystem.Get(), historyFile); err != nil {
		return err
	} else if !exists {
		if err = filesystem.Get().MkdirAll(filepath.Dir(historyFile), os.ModePerm); err != nil {
			return err
		}
	}

	err = afero.WriteFile(filesystem.Get(), historyFile, historyJSON, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

func ReadHistory() (map[string]*Entry, error) {
	historyFile, err := util.HistoryFilePath()

	if err != nil {
		return nil, err
	}

	// check if exists
	if exists, err := afero.Exists(filesystem.Get(), historyFile); err != nil {
		return nil, err
	} else if !exists {
		return make(map[string]*Entry), nil
	}

	history, err := afero.ReadFile(filesystem.Get(), historyFile)

	if err != nil {
		return nil, err
	}

	var historyEntries map[string]*Entry

	err = json.Unmarshal(history, &historyEntries)

	for k, entry := range historyEntries {
		s, ok := util.Find(config.UserConfig.Scrapers, func(s *scraper.Scraper) bool {
			return s.Source.Name == entry.SourceName
		})

		if !ok {
			delete(historyEntries, k)
		}

		entry.Manga.Scraper = s
		entry.Chapter.Scraper = s
	}

	if err != nil {
		return nil, err
	}

	return historyEntries, nil
}
