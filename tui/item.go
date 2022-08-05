package tui

import (
	"github.com/metafates/mangal/history"
	"github.com/metafates/mangal/source"
)

type listItem struct {
	title       string
	description string
	internal    interface{}
}

func (t *listItem) Title() string {
	if t.title != "" {
		return t.title
	}

	switch t.internal.(type) {
	case source.Source:
		return t.internal.(source.Source).Name()
	case *source.Manga:
		return t.internal.(*source.Manga).Name
	case *source.Chapter:
		return t.internal.(*source.Chapter).Name
	case *history.SavedChapter:
		return t.internal.(*history.SavedChapter).MangaName
	}

	panic(t.internal)
}

func (t *listItem) Description() string {
	if t.description != "" {
		return t.description
	}

	switch t.internal.(type) {
	case source.Source:
		return t.internal.(source.Source).ID()
	case *source.Manga:
		return t.internal.(*source.Manga).URL
	case *source.Chapter:
		return t.internal.(*source.Chapter).URL
	case *history.SavedChapter:
		return t.internal.(*history.SavedChapter).Name
	}

	panic("unreachable")
}

func (t *listItem) FilterValue() string {
	return t.Title()
}
