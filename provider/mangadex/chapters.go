package mangadex

import (
	"fmt"
	"github.com/darylhjd/mangodex"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/source"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"net/url"
	"strconv"
)

func (m *Mangadex) ChaptersOf(manga *source.Manga) ([]*source.Chapter, error) {
	if cached, ok := m.cachedChapters[manga.URL]; ok {
		return cached, nil
	}

	params := url.Values{}
	params.Set("limit", strconv.Itoa(500))
	ratings := []string{mangodex.Safe, mangodex.Suggestive}
	for _, rating := range ratings {
		params.Add("contentRating[]", rating)
	}

	if viper.GetBool(constant.MangadexNSFW) {
		params.Add("contentRating[]", mangodex.Porn)
		params.Add("contentRating[]", mangodex.Erotica)
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
			// Skip external chapters. Their pages cannot be downloaded.
			if chapter.Attributes.ExternalURL != nil && !viper.GetBool(constant.MangadexShowUnavailableChapters) {
				continue
			}

			// skip chapters that are not in the current language
			if chapter.Attributes.TranslatedLanguage != viper.GetString(constant.MangadexLanguage) {
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
				Index:    n + 1,
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

	m.cachedChapters[manga.URL] = chapters
	return chapters, nil
}
