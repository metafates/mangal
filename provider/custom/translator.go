package custom

import (
	"errors"
	"github.com/metafates/mangal/source"
	lua "github.com/yuin/gopher-lua"
	"strconv"
	"strings"
)

func mangaFromTable(table *lua.LTable, index uint16) (*source.Manga, error) {
	name := table.RawGetString("name")

	if name.Type() != lua.LTString {
		return nil, errors.New("type of field \"name\" should be string")
	}

	url := table.RawGetString("url")
	if url.Type() != lua.LTString {
		return nil, errors.New("type of field \"url\" should be string")
	}

	return &source.Manga{
		Name:     strings.TrimSpace(name.String()),
		URL:      strings.TrimSpace(url.String()),
		Index:    index,
		Chapters: []*source.Chapter{},
	}, nil
}

func chapterFromTable(table *lua.LTable, manga *source.Manga, index uint16) (*source.Chapter, error) {
	name := table.RawGetString("name")

	if name.Type() != lua.LTString {
		return nil, errors.New("type of field \"name\" should be string")
	}

	url := table.RawGetString("url")
	if url.Type() != lua.LTString {
		return nil, errors.New("type of field \"url\" should be string")
	}

	chapter := &source.Chapter{
		Name:  strings.TrimSpace(name.String()),
		URL:   strings.TrimSpace(url.String()),
		Manga: manga,
		Index: index,
		Pages: []*source.Page{},
	}

	manga.Chapters = append(manga.Chapters, chapter)
	return chapter, nil
}

func pageFromTable(table *lua.LTable, chapter *source.Chapter) (*source.Page, error) {
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

	page := &source.Page{
		URL:       strings.TrimSpace(url.String()),
		Index:     uint16(num),
		Chapter:   chapter,
		Extension: ".jpg",
	}

	chapter.Pages = append(chapter.Pages, page)
	return page, nil
}
