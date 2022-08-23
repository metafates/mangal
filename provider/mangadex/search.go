package mangadex

import (
	"fmt"
	"github.com/darylhjd/mangodex"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/source"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"strconv"
)

func (m *Mangadex) Search(query string) ([]*source.Manga, error) {
	if cached, ok := m.cachedMangas[query]; ok {
		return cached, nil
	}

	params := url.Values{}
	params.Set("limit", strconv.Itoa(100))

	ratings := []string{mangodex.Safe, mangodex.Suggestive}

	for _, rating := range ratings {
		params.Add("contentRating[]", rating)
	}

	if viper.GetBool(constant.MangadexNSFW) {
		params.Add("contentRating[]", mangodex.Porn)
		params.Add("contentRating[]", mangodex.Erotica)
	}

	params.Set("order[followedCount]", "desc")
	params.Set("title", query)

	mangaList, err := m.client.Manga.GetMangaList(params)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	var mangas []*source.Manga

	for i, manga := range mangaList.Data {
		m := source.Manga{
			Name:     manga.GetTitle(viper.GetString(constant.MangadexLanguage)),
			URL:      fmt.Sprintf("https://mangadex.org/title/%s", manga.ID),
			Index:    uint16(i),
			SourceID: ID,
			ID:       manga.ID,
		}

		mangas = append(mangas, &m)
	}

	m.cachedMangas[query] = mangas
	return mangas, nil
}
