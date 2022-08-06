package mangadex

import (
	"fmt"
	"github.com/darylhjd/mangodex"
	"github.com/metafates/mangal/source"
	"golang.org/x/exp/slices"
	"net/url"
	"strconv"
)

func (m *Mangadex) ChaptersOf(manga *source.Manga) ([]*source.Chapter, error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(500))
	ratings := []string{mangodex.Safe, mangodex.Suggestive, mangodex.Erotica}
	for _, rating := range ratings {
		params.Add("contentRating[]", rating)
	}

	// scanlation group for the chapter
	params.Add("includes[]", mangodex.ScanlationGroupRel)
	params.Set("order[chapter]", "asc")

	var chapters []*source.Chapter
	var currOffset = 0

	for {
		params.Set("offset", strconv.Itoa(currOffset))
		list, err := m.client.Chapter.GetMangaChapters(manga.ID, params)
		if err != nil {
			return nil, err
		}

		for i, chapter := range list.Data {
			if chapter.Attributes.TranslatedLanguage != "en" {
				currOffset += 500
				continue
			}

			num, err := strconv.ParseUint(chapter.GetChapterNum(), 10, 16)
			n := uint16(num)
			if err != nil {
				n = uint16(i)
			}

			name := chapter.GetTitle()
			if name == "" {
				name = fmt.Sprintf("Chapter %d", n)
			} else {
				name = fmt.Sprintf("Chapter %d - %s", n, name)
			}

			chapters = append(chapters, &source.Chapter{
				Name:     name,
				Index:    n,
				SourceID: ID,
				ID:       chapter.ID,
				URL:      fmt.Sprintf("https://mangadex.org/chapter/%s", chapter.ID),
				Manga:    manga,
			})
		}
		currOffset += 500
		if currOffset >= list.Total {
			break
		}

		if currOffset >= list.Total {
			break
		}
	}

	slices.SortFunc(chapters, func(a, b *source.Chapter) bool {
		return a.Index < b.Index
	})

	return chapters, nil
}
