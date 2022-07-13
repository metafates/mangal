package tui

import (
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/scraper"
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/util"
	"strconv"
	"strings"
)

// listItem is a list item used in the manga and chapters lists
// It contains the URL of the manga/chapter and the title of the manga/chapter
type listItem struct {
	selected bool
	url      *scraper.URL
}

func (l *listItem) Select() {
	l.selected = !l.selected
}
func (l *listItem) Title() string {
	var title string

	if l.selected {
		title = style.Accent.Bold(true).Render(config.UserConfig.UI.Mark) + " " + l.url.Info
	} else {
		title = l.url.Info
	}

	// If it's a manga
	if l.url.Relation == nil {
		return title
	}

	index := l.url.Index

	// replace according to the name template
	template := strings.ReplaceAll(config.UserConfig.UI.ChapterNameTemplate, "%0d", util.PadZeros(index, 4))
	template = strings.ReplaceAll(template, "%d", strconv.Itoa(index))
	template = strings.ReplaceAll(template, "%s", title)

	return template
}

func (l *listItem) Description() string {
	return l.url.Address
}

func (l *listItem) FilterValue() string {
	return l.url.Info
}
