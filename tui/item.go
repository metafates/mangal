package tui

import (
	"fmt"
	"github.com/metafates/mangal/anilist"
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/style"
)

type listItem struct {
	title       string
	description string
	internal    interface{}
	marked      bool
}

func (t *listItem) toggleMark() {
	t.marked = !t.marked
}

func (t *listItem) Title() (title string) {
	switch e := t.internal.(type) {
	case *source.Chapter:
		title = fmt.Sprintf("%s %s", e.Name, style.Faint(e.Volume))
	case *source.Manga:
		title = e.Name
	case *history.SavedChapter:
		title = e.MangaName
	case *anilist.Manga:
		title = e.Title.English
	default:
		title = t.title
	}

	if title != "" && t.marked {
		title = fmt.Sprintf("%s %s", title, icon.Get(icon.Mark))
	}

	return
}

func (t *listItem) Description() string {
	switch e := t.internal.(type) {
	case *source.Chapter:
		return e.URL
	case *source.Manga:
		return e.URL
	case *history.SavedChapter:
		return fmt.Sprintf("%s : %d / %d", e.Name, e.Index, e.MangaChaptersTotal)
	case *anilist.Manga:
		return e.SiteURL
	default:
		return t.description
	}
}

func (t *listItem) FilterValue() string {
	return t.Title()
}
