package manganelo

import (
	"fmt"
	"github.com/metafates/mangal/source"
	"net/url"
	"strings"
)

func (m *Manganelo) Search(query string) ([]*source.Manga, error) {
	address := fmt.Sprintf("https://ww5.manganelo.tv/search/%s", url.QueryEscape(strings.TrimSpace(strings.ToLower(query))))

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
