package source

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
)

type Chapter struct {
	Name  string
	URL   string
	Manga *Manga
}

func chapterFromTable(table *lua.LTable, mangaRelation *Manga) (*Chapter, error) {
	name := table.RawGetString("name")

	if name.Type() != lua.LTString {
		return nil, errors.New("type of field \"name\" should be string")
	}

	url := table.RawGetString("url")
	if url.Type() != lua.LTString {
		return nil, errors.New("type of field \"url\" should be string")
	}

	return &Chapter{
		Name:  name.String(),
		URL:   url.String(),
		Manga: mangaRelation,
	}, nil
}
