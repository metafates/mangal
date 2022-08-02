package source

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
)

type Manga struct {
	Name string
	URL  string
}

func mangaFromTable(table *lua.LTable) (*Manga, error) {
	name := table.RawGetString("name")

	if name.Type() != lua.LTString {
		return nil, errors.New("type of field \"name\" should be string")
	}

	url := table.RawGetString("url")
	if url.Type() != lua.LTString {
		return nil, errors.New("type of field \"url\" should be string")
	}

	return &Manga{
		Name: name.String(),
		URL:  url.String(),
	}, nil
}
