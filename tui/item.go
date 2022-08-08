package tui

import "github.com/metafates/mangal/icon"

type listItem struct {
	title       string
	description string
	internal    interface{}
	marked      bool
}

func (t *listItem) toggleMark() {
	t.marked = !t.marked
}

func (t *listItem) Title() string {
	if t.title != "" {
		if t.marked {
			return icon.Get(icon.Mark) + " " + t.title
		} else {
			return t.title
		}
	}

	panic("title is empty")
}

func (t *listItem) Description() string {
	if t.description != "" {
		return t.description
	}

	panic("description is empty")
}

func (t *listItem) FilterValue() string {
	return t.Title()
}
