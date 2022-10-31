package generic

import (
	"github.com/metafates/mangal/source"
)

// Search for mangas by given title
func (s *Scraper) Search(query string) ([]*source.Manga, error) {
	address := s.config.GenerateSearchURL(query)

	if urls, ok := s.mangas[address]; ok {
		return urls, nil
	}

	err := s.mangasCollector.Visit(address)

	if err != nil {
		return nil, err
	}

	s.mangasCollector.Wait()
	return s.mangas[address], nil
}
