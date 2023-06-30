package anilistmangas

import (
	"fmt"
	"github.com/mangalorg/libmangal"
)

type Item struct {
	Manga *libmangal.AnilistManga
}

func (i Item) FilterValue() string {
	for _, title := range []string{
		i.Manga.Title.English,
		i.Manga.Title.Romaji,
		i.Manga.Title.Native,
	} {
		if title != "" {
			return title
		}
	}

	return "Untitled"
}

func (i Item) Title() string {
	return i.FilterValue()
}

func (i Item) Description() string {
	return fmt.Sprint("https://anilist.co/manga/", i.Manga.ID)
}
