package source

// Source is the interface that all sources must implement.
type Source interface {
	Name() string
	Search(query string) ([]*Manga, error)
	ChaptersOf(manga *Manga) ([]*Chapter, error)
	PagesOf(chapter *Chapter) ([]*Page, error)
	ID() string
}
