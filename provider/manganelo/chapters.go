package manganelo

import (
	"github.com/gocolly/colly"
	"github.com/metafates/mangal/source"
)

func (m *Manganelo) ChaptersOf(manga *source.Manga) ([]*source.Chapter, error) {
	if chapters, ok := m.chapters[manga.URL]; ok {
		return chapters, nil
	}

	ctx := colly.NewContext()
	ctx.Put("manga", manga)
	err := m.chaptersCollector.Request("GET", manga.URL, nil, ctx, nil)

	if err != nil {
		return nil, err
	}

	m.chaptersCollector.Wait()

	// reverse chapters
	// will happend only once for each manga
	chapters := m.chapters[manga.URL]
	reversed := make([]*source.Chapter, len(chapters))
	for i, chapter := range chapters {
		reversed[len(chapters)-i-1] = chapter
		chapter.Index = uint16(len(chapters) - i - 1)
	}

	return reversed, nil
}
