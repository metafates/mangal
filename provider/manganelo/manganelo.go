package manganelo

import "github.com/metafates/mangal/source"

type Manganelo struct{}

func New() source.Source {
	return &Manganelo{}
}

func (_ Manganelo) Name() string {
	return "manganelo"
}

func (m Manganelo) Search(query string) ([]*source.Manga, error) {
	return nil, nil
}

func (m Manganelo) ChaptersOf(manga *source.Manga) ([]*source.Chapter, error) {
	return nil, nil
}

func (m Manganelo) PagesOf(chapter *source.Chapter) ([]*source.Page, error) {
	return nil, nil
}
