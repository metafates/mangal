package main

import "encoding/json"

type HistoryEntry struct {
	Manga   *URL
	Chapter int
}

func WriteHistory(manga *URL, chapter *URL) error {
	historyFile, err := HistoryFile()

	if err != nil {
		return err
	}

	history, err := ReadHistory()

	if err != nil {
		return err
	}

	history[manga.Address] = &HistoryEntry{
		Manga:   manga,
		Chapter: chapter.Index,
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

	if err != nil {
		return nil, err
	}

	return historyEntries, nil
}
