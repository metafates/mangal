package mangadex

import (
	"bytes"
	"errors"
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

	if len(downloader.Pages) == 0 {
		return nil, errors.New("there were no pages for this chapter")
	}

	for i, name := range downloader.Pages {
		image, err := downloader.GetChapterPage(name)
		if err != nil {
			return nil, err
		}

		if len(image) == 0 {
			return nil, errors.New("image is empty")
		}

		page := source.Page{
			Index:     uint16(i),
			Chapter:   chapter,
			Extension: filepath.Ext(name),
			Contents:  io.NopCloser(bytes.NewReader(image)),
			Size:      uint64(len(image)),
			SourceID:  ID,
		}
		chapter.Pages = append(chapter.Pages, &page)
		pages = append(pages, &page)
	}

	return pages, nil
}
