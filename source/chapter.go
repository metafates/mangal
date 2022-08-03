package source

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
)

type Chapter struct {
	Name  string
	URL   string
	Manga *Manga
	Pages []*Page
}

func chapterFromTable(table *lua.LTable, manga *Manga) (*Chapter, error) {
	name := table.RawGetString("name")

	if name.Type() != lua.LTString {
		return nil, errors.New("type of field \"name\" should be string")
	}

	url := table.RawGetString("url")
	if url.Type() != lua.LTString {
		return nil, errors.New("type of field \"url\" should be string")
	}

	chapter := &Chapter{
		Name:  name.String(),
		URL:   url.String(),
		Manga: manga,
		Pages: []*Page{},
	}

	manga.Chapters = append(manga.Chapters, chapter)
	return chapter, nil
}
