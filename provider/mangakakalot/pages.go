package mangakakalot

import (
	"github.com/gocolly/colly"
	"github.com/metafates/mangal/source"
	"net/http"
)

func (m *Mangakakalot) PagesOf(chapter *source.Chapter) ([]*source.Page, error) {
	if pages, ok := m.pages[chapter.URL]; ok {
		return pages, nil
	}

	ctx := colly.NewContext()
	ctx.Put("chapter", chapter)
	err := m.pagesCollector.Request(http.MethodGet, chapter.URL, nil, ctx, nil)

	if err != nil {
		return nil, err
	}

	m.pagesCollector.Wait()

	return m.pages[chapter.URL], nil
}
