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
		Romaji string `json:"romaji" jsonschema:"description=Romanized title of the manga."`
		// English is the english title of the manga.
		English string `json:"english" jsonschema:"description=English title of the manga."`
		// Native is the native title of the manga. (Usually in kanji)
		Native string `json:"native" jsonschema:"description=Native title of the manga. Usually in kanji."`
	} `json:"title"`
	// ID is the id of the manga on Anilist.
	ID int `json:"id" jsonschema:"description=ID of the manga on Anilist."`
	// Description is the description of the manga in html format.
	Description string `json:"description" jsonschema:"description=Description of the manga in html format."`
	// CoverImage is the cover image of the manga.
	CoverImage struct {
		// ExtraLarge is the url of the extra large cover image.
		// If the image is not available, large will be used instead.
		ExtraLarge string `json:"extraLarge" jsonschema:"description=URL of the extra large cover image. If the image is not available, large will be used instead."`
		// Large is the url of the large cover image.
		Large string `json:"large" jsonschema:"description=URL of the large cover image."`
		// Medium is the url of the medium cover image.
		Medium string `json:"medium" jsonschema:"description=URL of the medium cover image."`
		// Color is the average color of the cover image.
		Color string `json:"color" jsonschema:"description=Average color of the cover image."`
	} `json:"coverImage" jsonschema:"description=Cover image of the manga."`
	// BannerImage of the media
	BannerImage string `json:"bannerImage" jsonschema:"description=Banner image of the manga."`
	// Tags are the tags of the manga.
	Tags []struct {
		// Name of the tag.
		Name string `json:"name" jsonschema:"description=Name of the tag."`
		// Description of the tag.
		Description string `json:"description" jsonschema:"description=Description of the tag."`
		// Rank of the tag. How relevant it is to the manga from 1 to 100.
		Rank int `json:"rank" jsonschema:"description=Rank of the tag. How relevant it is to the manga from 1 to 100."`
	} `json:"tags"`
	// Genres of the manga
	Genres []string `json:"genres" jsonschema:"description=Genres of the manga."`
	// Characters are the primary characters of the manga.
	Characters struct {
		Nodes []struct {
			Name struct {
				// Full is the full name of the character.
				Full string `json:"full" jsonschema:"description=Full name of the character."`
				// Native is the native name of the character. Usually in kanji.
				Native string `json:"native" jsonschema:"description=Native name of the character. Usually in kanji."`
			} `json:"name"`
		} `json:"nodes"`
	} `json:"characters"`
	Staff struct {
		Edges []struct {
			Role string `json:"role" jsonschema:"description=Role of the staff member."`
			Node struct {
				Name struct {
					Full string `json:"full" jsonschema:"description=Full name of the staff member."`
				} `json:"name"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"staff"`
	// StartDate is the date the manga started publishing.
	StartDate date `json:"startDate" jsonschema:"description=Date the manga started publishing."`
	// EndDate is the date the manga ended publishing.
	EndDate date `json:"endDate" jsonschema:"description=Date the manga ended publishing."`
	// Synonyms are the synonyms of the manga (Alternative titles).
	Synonyms []string `json:"synonyms" jsonschema:"description=Synonyms of the manga (Alternative titles)."`
	// Status is the status of the manga. (FINISHED, RELEASING, NOT_YET_RELEASED, CANCELLED)
	Status string `json:"status" jsonschema:"enum=FINISHED,enum=RELEASING,enum=NOT_YET_RELEASED,enum=CANCELLED,enum=HIATUS"`
	// IDMal is the id of the manga on MyAnimeList.
	IDMal int `json:"idMal" jsonschema:"description=ID of the manga on MyAnimeList."`
	// Chapters is the amount of chapters the manga has when complete.
	Chapters int `json:"chapters" jsonschema:"description=Amount of chapters the manga has when complete."`
	// SiteURL is the url of the manga on Anilist.
	SiteURL string `json:"siteUrl" jsonschema:"description=URL of the manga on Anilist."`
	// Country of origin of the manga.
	Country string `json:"countryOfOrigin" jsonschema:"description=Country of origin of the manga."`
	// External urls related to the manga.
	External []struct {
		URL string `json:"url" jsonschema:"description=URL of the external link."`
	} `json:"externalLinks" jsonschema:"description=External links related to the manga."`
}

// Name of the manga. If English is available, it will be used. Otherwise, Romaji will be used.
func (m *Manga) Name() string {
	if m.Title.English == "" {
		return m.Title.Romaji
	}

	return m.Title.English
}
