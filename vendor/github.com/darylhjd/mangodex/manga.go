package mangodex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	MangaListPath            = "manga"
	CheckIfMangaFollowedPath = "user/follows/manga/%s"
	ToggleMangaFollowPath    = "manga/%s/follow"
)

// MangaService : Provides Manga services provided by the API.
type MangaService service

// MangaList : A response for getting a list of manga.
type MangaList struct {
	Result   string  `json:"result"`
	Response string  `json:"response"`
	Data     []Manga `json:"data"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Total    int     `json:"total"`
}

func (ml *MangaList) GetResult() string {
	return ml.Result
}

// Manga : Struct containing information on a Manga.
type Manga struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Attributes    MangaAttributes `json:"attributes"`
	Relationships []Relationship  `json:"relationships"`
}

// GetTitle : Get title of the Manga.
func (m *Manga) GetTitle(langCode string) string {
	if title := m.Attributes.Title.GetLocalString(langCode); title != "" {
		return title
	}
	return m.Attributes.AltTitles.GetLocalString(langCode)
}

// GetDescription : Get description of the Manga.
func (m *Manga) GetDescription(langCode string) string {
	return m.Attributes.Description.GetLocalString(langCode)
}

// MangaAttributes : Attributes for a Manga.
type MangaAttributes struct {
	Title                  LocalisedStrings `json:"title"`
	AltTitles              LocalisedStrings `json:"altTitles"`
	Description            LocalisedStrings `json:"description"`
	IsLocked               bool             `json:"isLocked"`
	Links                  LocalisedStrings `json:"links"`
	OriginalLanguage       string           `json:"originalLanguage"`
	LastVolume             *string          `json:"lastVolume"`
	LastChapter            *string          `json:"lastChapter"`
	PublicationDemographic *string          `json:"publicationDemographic"`
	Status                 *string          `json:"status"`
	Year                   *int             `json:"year"`
	ContentRating          *string          `json:"contentRating"`
	Tags                   []Tag            `json:"tags"`
	State                  string           `json:"state"`
	Version                int              `json:"version"`
	CreatedAt              string           `json:"createdAt"`
	UpdatedAt              string           `json:"updatedAt"`
}

// GetMangaList : Get a list of Manga.
// https://api.mangadex.org/docs.html#operation/get-search-manga
func (s *MangaService) GetMangaList(params url.Values) (*MangaList, error) {
	return s.GetMangaListContext(context.Background(), params)
}

// GetMangaListContext : GetMangaList with custom context.
func (s *MangaService) GetMangaListContext(ctx context.Context, params url.Values) (*MangaList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = MangaListPath

	// Set query parameters
	u.RawQuery = params.Encode()

	var l MangaList
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// CheckIfMangaFollowed : Check if a user follows a manga.
func (s *MangaService) CheckIfMangaFollowed(id string) (bool, error) {
	return s.CheckIfMangaFollowedContext(context.Background(), id)
}

// CheckIfMangaFollowedContext : CheckIfMangaFollowed with custom context.
func (s *MangaService) CheckIfMangaFollowedContext(ctx context.Context, id string) (bool, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(CheckIfMangaFollowedPath, id)

	var r Response
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ToggleMangaFollowStatus :Toggle follow status for a manga.
func (s *MangaService) ToggleMangaFollowStatus(id string, toFollow bool) (*Response, error) {
	return s.ToggleMangaFollowStatusContext(context.Background(), id, toFollow)
}

// ToggleMangaFollowStatusContext  ToggleMangaFollowStatus with custom context.
func (s *MangaService) ToggleMangaFollowStatusContext(ctx context.Context, id string, toFollow bool) (*Response, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(ToggleMangaFollowPath, id)

	method := http.MethodPost // To follow
	if !toFollow {
		method = http.MethodDelete // To unfollow
	}

	var r Response
	err := s.client.RequestAndDecode(ctx, method, u.String(), nil, &r)
	return &r, err
}
