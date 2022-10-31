package custom

import (
	"fmt"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

type mapping lo.Tuple4[lua.LValueType, bool, func(string) error, string]

func translate(
	table *lua.LTable,
	mappings map[string]mapping,
) (err error) {
	for field, t := range mappings {
		var (
			type_    = t.A
			required = t.B
			handle   = t.C
			default_ = t.D
		)

		val := table.RawGetString(field)
		if val.Type() == lua.LTNil {
			if required {
				err = fmt.Errorf(`field of "%s" is required`, field)
			} else {
				err = handle(default_)
			}
		} else if val.Type() != type_ {
			err = fmt.Errorf(`field of "%s" must be of type %s`, field, type_)
		} else {
			err = handle(val.String())
		}

		if err != nil {
			return
		}
	}

	return
}

func mangaFromTable(table *lua.LTable, index uint16) (manga *source.Manga, err error) {
	manga = &source.Manga{
		Index:    index,
		Chapters: []*source.Chapter{},
	}

	mappings := map[string]mapping{
		"name":    {A: lua.LTString, B: true, C: func(v string) error { manga.Name = v; return nil }},
		"url":     {A: lua.LTString, B: true, C: func(v string) error { manga.URL = v; return nil }},
		"summary": {A: lua.LTString, B: false, C: func(v string) error { manga.Metadata.Summary = v; return nil }},
		"cover": {A: lua.LTString, B: false, C: func(v string) error {
			if v == "" {
				return nil
			}

			_, err := url.Parse(v)
			if err != nil {
				return err
			}

			manga.Metadata.Cover.ExtraLarge = v
			return nil
		}},
		"genres": {A: lua.LTString, B: false, C: func(v string) error {
			manga.Metadata.Genres = lo.Map(strings.Split(v, ","), func(genre string, _ int) string {
				return strings.TrimSpace(genre)
			})
			return nil
		}},
	}

	err = translate(table, mappings)
	return
}

func chapterFromTable(table *lua.LTable, manga *source.Manga, index uint16) (chapter *source.Chapter, err error) {
	chapter = &source.Chapter{
		Manga: manga,
		Index: index,
		Pages: []*source.Page{},
	}

	mappings := map[string]mapping{
		"name":          {A: lua.LTString, B: true, C: func(v string) error { chapter.Name = v; return nil }},
		"url":           {A: lua.LTString, B: true, C: func(v string) error { chapter.URL = v; return nil }},
		"volume":        {A: lua.LTString, B: false, C: func(v string) error { chapter.Volume = v; return nil }},
		"manga_summary": {A: lua.LTString, B: false, C: func(v string) error { manga.Metadata.Summary = v; return nil }},
		"manga_genres": {A: lua.LTString, B: false, C: func(v string) error {
			manga.Metadata.Genres = lo.Map(strings.Split(v, ","), func(genre string, _ int) string {
				return strings.TrimSpace(genre)
			})
			return nil
		}},
		"manga_cover": {A: lua.LTString, B: false, C: func(v string) error {
			if v == "" {
				return nil
			}
			_, err := url.Parse(v)
			if err != nil {
				return err
			}

			manga.Metadata.Cover.ExtraLarge = v
			return nil
		}},
	}

	err = translate(table, mappings)
	manga.Chapters = append(manga.Chapters, chapter)
	return
}

func pageFromTable(table *lua.LTable, chapter *source.Chapter) (page *source.Page, err error) {
	page = &source.Page{
		Chapter: chapter,
	}

	mappings := map[string]mapping{
		"url": {A: lua.LTString, B: true, C: func(v string) error { page.URL = v; return nil }},
		"index": {A: lua.LTNumber, B: true, C: func(v string) error {
			num, err := strconv.ParseUint(v, 10, 16)
			if err != nil {
				return err
			}

			page.Index = uint16(num)
			return nil
		}},
	}

	err = translate(table, mappings)
	if err != nil {
		return
	}

	page.Extension = filepath.Ext(page.URL)
	chapter.Pages = append(chapter.Pages, page)
	return
}
