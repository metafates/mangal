package anilist

type Manga struct {
	URL   string `json:"url"`
	Title struct {
		Romaji  string `json:"romaji"`
		English string `json:"english"`
		Native  string `json:"native"`
	} `json:"title"`
	ID          int    `json:"id"`
	Description string `json:"description"`
	CoverImage  struct {
		ExtraLarge string `json:"extraLarge"`
	} `json:"coverImage"`
	Tags []struct {
		Name string `json:"name"`
	} `json:"tags"`
	Genres     []string `json:"genres"`
	Characters struct {
		Nodes []struct {
			Name struct {
				Full string `json:"full"`
			} `json:"name"`
		} `json:"nodes"`
	} `json:"characters"`
	StartDate struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
	} `json:"startDate"`
	Status string `json:"status"`
}

func (m *Manga) Name() string {
	if m.Title.English == "" {
		return m.Title.Romaji
	}

	return m.Title.English
}
