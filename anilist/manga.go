package anilist

type date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type Manga struct {
	// Title of the manga
	Title struct {
		// Romaji is the romanized title of the manga.
		Romaji string `json:"romaji"`
		// English is the english title of the manga.
		English string `json:"english"`
		// Native is the native title of the manga. (Usually in kanji)
		Native string `json:"native"`
	} `json:"title"`
	// ID is the id of the manga on Anilist.
	ID int `json:"id"`
	// Description is the description of the manga in html format.
	Description string `json:"description"`
	// CoverImage is the cover image of the manga.
	CoverImage struct {
		// ExtraLarge is the url of the extra large cover image.
		// If the image is not available, large will be used instead.
		ExtraLarge string `json:"extraLarge"`
		// Large is the url of the large cover image.
		Large string `json:"large"`
		// Medium is the url of the medium cover image.
		Medium string `json:"medium"`
		// Color is the average color of the cover image.
		Color string `json:"color"`
	} `json:"coverImage"`
	// BannerImage of the media
	BannerImage string `json:"bannerImage"`
	// Tags are the tags of the manga.
	Tags []struct {
		// Name of the tag.
		Name string `json:"name"`
		// Description of the tag.
		Description string `json:"description"`
		// Rank of the tag. How relevant it is to the manga from 1 to 100.
		Rank int `json:"rank"`
	} `json:"tags"`
	// Genres of the manga
	Genres []string `json:"genres"`
	// Characters are the primary characters of the manga.
	Characters struct {
		Nodes []struct {
			Name struct {
				// Full is the full name of the character.
				Full string `json:"full"`
				// Native is the native name of the character. Usually in kanji.
				Native string `json:"native"`
			} `json:"name"`
		} `json:"nodes"`
	} `json:"characters"`
	Staff struct {
		Edges []struct {
			Role string `json:"role"`
			Node struct {
				Name struct {
					Full string `json:"full"`
				} `json:"name"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"staff"`
	// StartDate is the date the manga started publishing.
	StartDate date `json:"startDate"`
	// EndDate is the date the manga ended publishing.
	EndDate date `json:"endDate"`
	// Synonyms are the synonyms of the manga (Alternative titles).
	Synonyms []string `json:"synonyms"`
	// Status is the status of the manga. (FINISHED, RELEASING, NOT_YET_RELEASED, CANCELLED)
	Status string `json:"status" jsonschema:"enum=FINISHED,enum=RELEASING,enum=NOT_YET_RELEASED,enum=CANCELLED,enum=HIATUS"`
	// IDMal is the id of the manga on MyAnimeList.
	IDMal int `json:"idMal"`
	// Chapters is the amount of chapters the manga has when complete.
	Chapters int `json:"chapters"`
	// SiteURL is the url of the manga on Anilist.
	SiteURL string `json:"siteUrl"`
	// Country of origin of the manga.
	Country string `json:"countryOfOrigin"`
	// External urls related to the manga.
	External []struct {
		URL string `json:"url"`
	} `json:"externalLinks"`
}

// Name of the manga. If English is available, it will be used. Otherwise, Romaji will be used.
func (m *Manga) Name() string {
	if m.Title.English == "" {
		return m.Title.Romaji
	}

	return m.Title.English
}
