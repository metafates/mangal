package source

type Manga struct {
	Name     string
	URL      string
	Index    uint16
	SourceID string
	ID       string
	Chapters []*Chapter
}

func (m *Manga) String() string {
	return m.Name
}
