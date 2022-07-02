package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

type HistoryEntry struct {
	Manga      *URL
	Chapter    *URL
	SourceName string
}

func (h *HistoryEntry) Title() string {
	return h.Manga.Info
}

func (h *HistoryEntry) Description() string {
	// replace according to the name description
	description := strings.ReplaceAll(UserConfig.UI.ChapterNameTemplate, "%0d", PadZeros(h.Chapter.Index, 4))
	description = strings.ReplaceAll(description, "%d", strconv.Itoa(h.Chapter.Index))
	description = strings.ReplaceAll(description, "%s", h.Chapter.Info)

	return description
}

func (h *HistoryEntry) FilterValue() string {
	return h.Manga.Info
}

func WriteHistory(chapter *URL) error {
	historyFile, err := HistoryFile()

	if err != nil {
		return err
	}

	history, err := ReadHistory()

	if err != nil {
		return err
	}

	history[chapter.Relation.Address] = &HistoryEntry{
		Manga:      chapter.Relation,
		Chapter:    chapter,
		SourceName: chapter.Scraper.Source.Name,
	}

	historyJSON, err := json.Marshal(history)

	if err != nil {
		return err
	}

	err = Afero.WriteFile(historyFile, historyJSON, 0777)

	if err != nil {
		return err
	}

	return nil
}

func ReadHistory() (map[string]*HistoryEntry, error) {
	historyFile, err := HistoryFile()

	if err != nil {
		return nil, err
	}

	// check if exists
	if exists, err := Afero.Exists(historyFile); err != nil {
		return nil, err
	} else if !exists {
		return make(map[string]*HistoryEntry), nil
	}

	history, err := Afero.ReadFile(historyFile)

	if err != nil {
		return nil, err
	}

	var historyEntries map[string]*HistoryEntry

	err = json.Unmarshal(history, &historyEntries)

	for k, entry := range historyEntries {
		scraper, ok := Find(UserConfig.Scrapers, func(scraper *Scraper) bool {
			return scraper.Source.Name == entry.SourceName
		})

		if !ok {
			delete(historyEntries, k)
		}

		entry.Manga.Scraper = scraper
		entry.Chapter.Scraper = scraper
	}

	if err != nil {
		return nil, err
	}

	return historyEntries, nil
}
