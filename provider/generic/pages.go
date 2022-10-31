package generic

import (
	"github.com/gocolly/colly/v2"
	"github.com/metafates/mangal/source"
	"net/http"
)

// PagesOf given source.Chapter
func (s *Scraper) PagesOf(chapter *source.Chapter) ([]*source.Page, error) {
	if pages, ok := s.pages[chapter.URL]; ok {
		return pages, nil
	}

	ctx := colly.NewContext()
	ctx.Put("chapter", chapter)
	err := s.pagesCollector.Request(http.MethodGet, chapter.URL, nil, ctx, nil)

	if err != nil {
		return nil, err
	}

	s.pagesCollector.Wait()

	return s.pages[chapter.URL], nil
}
