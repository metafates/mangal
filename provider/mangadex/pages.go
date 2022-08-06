package mangadex

import (
	"bytes"
	"github.com/metafates/mangal/source"
	"io"
	"path/filepath"
)

func (m *Mangadex) PagesOf(chapter *source.Chapter) ([]*source.Page, error) {
	downloader, err := m.client.AtHome.NewMDHomeClient(chapter.ID, "data", false)
	if err != nil {
		return nil, err
	}

	var pages []*source.Page

	for i, name := range downloader.Pages {
		image, err := downloader.GetChapterPage(name)
		if err != nil {
			return nil, err
		}
		page := source.Page{
			Index:     uint16(i),
			Chapter:   chapter,
			Extension: filepath.Ext(name),
			Contents:  io.NopCloser(bytes.NewReader(image)),
			SourceID:  ID,
		}
		chapter.Pages = append(chapter.Pages, &page)
		pages = append(pages, &page)
	}

	return pages, nil
}
