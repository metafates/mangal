package source

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
	"strings"
	"sync"
)

type Chapter struct {
	Name  string
	URL   string
	Index uint16
	Manga *Manga
	Pages []*Page
}

func chapterFromTable(table *lua.LTable, manga *Manga, index uint16) (*Chapter, error) {
	name := table.RawGetString("name")

	if name.Type() != lua.LTString {
		return nil, errors.New("type of field \"name\" should be string")
	}

	url := table.RawGetString("url")
	if url.Type() != lua.LTString {
		return nil, errors.New("type of field \"url\" should be string")
	}

	chapter := &Chapter{
		Name:  strings.TrimSpace(name.String()),
		URL:   strings.TrimSpace(url.String()),
		Manga: manga,
		Index: index,
		Pages: []*Page{},
	}

	manga.Chapters = append(manga.Chapters, chapter)
	return chapter, nil
}

// DownloadPages downloads the Pages contents of the Chapter.
// Pages needs to be set before calling this function.
func (c *Chapter) DownloadPages() error {
	wg := sync.WaitGroup{}
	wg.Add(len(c.Pages))

	var err error
	for _, page := range c.Pages {
		go func(page *Page) {
			defer wg.Done()

			// if at any point, an error is encountered, stop downloading other pages
			if err != nil {
				return
			}

			err = page.Download()
		}(page)
	}

	wg.Wait()
	return err
}
