package source

// Manga is a manga from a source.
type Manga struct {
	// Name of the manga
	Name string
	// URL of the manga
	URL string
	// Index of the manga in the source.
	Index uint16
	// SourceID of the source the manga is from.
	SourceID string
	// ID of manga in the source.
	ID string
	// Chapters of the manga
	Chapters []*Chapter
}

func (m *Manga) String() string {
	return m.Name
}
