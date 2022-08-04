package source

type Source interface {
	Name() string
	Search(query string) ([]*Manga, error)
	ChaptersOf(manga *Manga) ([]*Chapter, error)
	PagesOf(chapter *Chapter) ([]*Page, error)
	ID() string
}
