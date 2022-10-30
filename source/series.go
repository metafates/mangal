package source

type SeriesJSON struct {
	Metadata struct {
		Type                 string `json:"type"`
		Name                 string `json:"name"`
		DescriptionFormatted string `json:"description_formatted"`
		DescriptionText      string `json:"description_text"`
		Status               string `json:"status"`
		Year                 int    `json:"year"`
		ComicImage           string `json:"ComicImage"`
		Publisher            string `json:"publisher"`
		ComicID              int    `json:"comicId"`
		BookType             string `json:"booktype"`
		TotalIssues          int    `json:"total_issues"`
		PublicationRun       string `json:"publication_run"`
	} `json:"metadata"`
}
