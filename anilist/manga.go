package anilist

type date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

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
		Large      string `json:"large"`
		Medium     string `json:"medium"`
		Color      string `json:"color"`
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
	StartDate date     `json:"startDate"`
	EndDate   date     `json:"endDate"`
	Synonyms  []string `json:"synonyms"`
	Status    string   `json:"status"`
	IDMal     int      `json:"idMal"`
	SiteURL   string   `json:"siteUrl"`
	Country   string   `json:"countryOfOrigin"`
	External  []struct {
		URL string `json:"url"`
	} `json:"externalLinks"`
}

func (m *Manga) Name() string {
	if m.Title.English == "" {
		return m.Title.Romaji
	}

	return m.Title.English
}
