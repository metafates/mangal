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
		title = e.Name()
	default:
		title = t.title
	}

	if title != "" && t.marked {
		title = fmt.Sprintf("%s %s", title, icon.Get(icon.Mark))
	}

	return
}

func (t *listItem) Description() (description string) {
	switch e := t.internal.(type) {
	case *source.Chapter:
		description = e.URL
	case *source.Manga:
		description = e.URL
	case *history.SavedChapter:
		description = fmt.Sprintf("%s : %d / %d", e.Name, e.Index, e.MangaChaptersTotal)
	case *anilist.Manga:
		description = e.SiteURL
	default:
		description = t.description
	}

	return
}

func (t *listItem) FilterValue() string {
	return t.Title()
}
