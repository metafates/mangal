package history

import (
	"fmt"
	"github.com/metafates/mangal/source"
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

func (c *SavedChapter) encode() string {
	return fmt.Sprintf("%s (%s)", c.MangaName, c.SourceID)
}

func (c *SavedChapter) String() string {
	return fmt.Sprintf("%s : %d / %d", c.MangaName, c.Index, c.MangaChaptersTotal)
}

func newSavedChapter(chapter *source.Chapter) *SavedChapter {
	return &SavedChapter{
		SourceID:           chapter.Manga.Source.ID(),
		MangaName:          chapter.Manga.Name,
		MangaURL:           chapter.Manga.URL,
		Name:               chapter.Name,
		URL:                chapter.URL,
		ID:                 chapter.ID,
		MangaID:            chapter.Manga.ID,
		MangaChaptersTotal: len(chapter.Manga.Chapters),
		Index:              int(chapter.Index),
	}
}