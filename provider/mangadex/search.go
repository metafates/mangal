package mangadex

import (
	"fmt"
	"github.com/darylhjd/mangodex"
	"github.com/metafates/mangal/source"
	"log"
	"net/url"
	"strconv"
)

func (m *Mangadex) Search(query string) ([]*source.Manga, error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(100))

	ratings := []string{mangodex.Safe, mangodex.Suggestive, mangodex.Erotica}

	for _, rating := range ratings {
		params.Add("contentRating[]", rating)
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
			Name:     manga.GetTitle("en"),
			URL:      fmt.Sprintf("https://mangadex.org/title/%s", manga.ID),
			Index:    uint16(i),
			SourceID: ID,
			ID:       manga.ID,
		}

		mangas = append(mangas, &m)
	}

	return mangas, nil
}
