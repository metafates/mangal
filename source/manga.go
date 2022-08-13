package source

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

type Manga struct {
	Name     string
	URL      string
	Index    uint16
	SourceID string
	ID       string
	Chapters []*Chapter
}

func (m *Manga) String() string {
	return m.Name
}

func mangaFromTable(table *lua.LTable, index uint16) (*Manga, error) {
	name := table.RawGetString("name")

	if name.Type() != lua.LTString {
		return nil, errors.New("type of field \"name\" should be string")
	}

	url := table.RawGetString("url")
	if url.Type() != lua.LTString {
		return nil, errors.New("type of field \"url\" should be string")
	}

	return &Manga{
		Name:     strings.TrimSpace(name.String()),
		URL:      strings.TrimSpace(url.String()),
		Index:    index,
		Chapters: []*Chapter{},
	}, nil
}
