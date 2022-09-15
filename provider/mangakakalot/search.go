package mangakakalot

import (
	"fmt"
	"github.com/metafates/mangal/source"
	"net/url"
	"strings"
)

func (m *Mangakakalot) Search(query string) ([]*source.Manga, error) {
	query = strings.ReplaceAll(query, " ", "_")
	address := fmt.Sprintf("https://mangakakalot.com/search/story/%s", url.QueryEscape(strings.TrimSpace(strings.ToLower(query))))

	if urls, ok := m.mangas[address]; ok {
		return urls, nil
	}

	err := m.mangasCollector.Visit(address)

	if err != nil {
		return nil, err
	}

	m.mangasCollector.Wait()
	return m.mangas[address], nil
}
