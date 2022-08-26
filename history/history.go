package history

import (
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/integration"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/where"
	"github.com/spf13/viper"
	"os"
)

type SavedChapter struct {
	SourceID           string `json:"source_id"`
	MangaName          string `json:"manga_name"`
	MangaURL           string `json:"manga_url"`
	MangaChaptersTotal int    `json:"manga_chapters_total"`
	Name               string `json:"name"`
	URL                string `json:"url"`
	ID                 string `json:"id"`
	Index              int    `json:"index"`
	MangaID            string `json:"manga_id"`
}

func (c *SavedChapter) String() string {
	return fmt.Sprintf("%s : %d / %d", c.MangaName, c.Index, c.MangaChaptersTotal)
}

// Get returns all chapters from the history file
func Get() (chapters map[string]*SavedChapter, err error) {
	log.Info("Getting history location")
	historyFile := where.History()

	// decode json into slice of structs
	log.Info("Reading history file")
	contents, err := filesystem.Get().ReadFile(historyFile)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Decoding history from json")
	err = json.Unmarshal(contents, &chapters)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

// Save saves the chapter to the history file
func Save(chapter *source.Chapter) error {
	if viper.GetBool(constant.AnilistEnable) {
		defer func() {
			log.Info("Saving chapter to anilist")
			err := integration.Anilist.MarkRead(chapter)
			if err != nil {
				log.Error("Saving chapter to anilist failed: " + err.Error())
			}
		}()
	}

	log.Info("Saving chapter to history")

	historyFile := where.History()

	// decode json into slice of structs
	var chapters map[string]*SavedChapter
	log.Info("Reading history file")
	contents, err := filesystem.Get().ReadFile(historyFile)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("Decoding history from json")
	err = json.Unmarshal(contents, &chapters)
	if err != nil {
		log.Error(err)
		return err
	}

	jsonChapter := SavedChapter{
		SourceID:           chapter.SourceID,
		MangaName:          chapter.Manga.Name,
		MangaURL:           chapter.Manga.URL,
		Name:               chapter.Name,
		URL:                chapter.URL,
		ID:                 chapter.ID,
		MangaID:            chapter.Manga.ID,
		MangaChaptersTotal: len(chapter.Manga.Chapters),
		Index:              int(chapter.Index),
	}

	chapters[fmt.Sprintf("%s (%s)", chapter.Manga.Name, chapter.SourceID)] = &jsonChapter

	// encode json
	log.Info("Encoding history to json")
	encoded, err := json.Marshal(chapters)
	if err != nil {
		log.Error(err)
		return err
	}

	// write to file
	log.Info("Writing history to file")
	err = filesystem.Get().WriteFile(historyFile, encoded, os.ModePerm)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// Remove removes the chapter from the history file
func Remove(chapter *SavedChapter) error {
	log.Info("Removing chapter from history")

	historyFile := where.History()

	// decode json into slice of structs
	var chapters map[string]*SavedChapter
	log.Info("Reading history file")
	contents, err := filesystem.Get().ReadFile(historyFile)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("Decoding history from json")
	err = json.Unmarshal(contents, &chapters)
	if err != nil {
		log.Error(err)
		return err
	}

	delete(chapters, fmt.Sprintf("%s (%s)", chapter.MangaName, chapter.SourceID))

	// encode json
	log.Info("Encoding history to json")
	encoded, err := json.Marshal(chapters)
	if err != nil {
		log.Error(err)
		return err
	}

	// write to file
	log.Info("Writing history to file")
	err = filesystem.Get().WriteFile(historyFile, encoded, os.ModePerm)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
