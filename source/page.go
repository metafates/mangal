package source

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
	"strconv"
)

type Page struct {
	URL     string
	Index   uint16
	Chapter *Chapter
}

func pageFromTable(table *lua.LTable, chapter *Chapter) (*Page, error) {
	url := table.RawGetString("url")

	if url.Type() != lua.LTString {
		return nil, errors.New("type of field \"url\" should be string")
	}

	index := table.RawGetString("index")

	if index.Type() != lua.LTNumber {
		return nil, errors.New("type of field \"index\" should be number")
	}

	num, err := strconv.ParseUint(index.String(), 10, 16)

	if err != nil {
		return nil, errors.New("index must be an unsigned 16 bit integer")
	}

	return &Page{
		URL:     url.String(),
		Index:   uint16(num),
		Chapter: chapter,
	}, nil
}
