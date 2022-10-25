package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	MangaChaptersPath    = "manga/%s/feed"
	MangaReadMarkersPath = "manga/%s/read"
)

// ChapterService : Provides Chapter services provided by the API.
type ChapterService service

// ChapterList : A response for getting a list of chapters.
type ChapterList struct {
	Result   string    `json:"result"`
	Response string    `json:"response"`
	Data     []Chapter `json:"data"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
	Total    int       `json:"total"`
}

func (cl *ChapterList) GetResult() string {
	return cl.Result
}

// Chapter : Struct containing information on a manga.
type Chapter struct {
	ID            string            `json:"id"`
	Type          string            `json:"type"`
	Attributes    ChapterAttributes `json:"attributes"`
	Relationships []Relationship    `json:"relationships"`
}

// GetTitle : Get a title for the chapter.
func (c *Chapter) GetTitle() string {
	return c.Attributes.Title
}

// GetChapterNum : Get the chapter's chapter number.
func (c *Chapter) GetChapterNum() string {
	if num := c.Attributes.Chapter; num != nil {
		return *num
	}
	return "-"
}

// ChapterAttributes : Attributes for a Chapter.
type ChapterAttributes struct {
	Title              string  `json:"title"`
	Volume             *string `json:"volume"`
	Chapter            *string `json:"chapter"`
	TranslatedLanguage string  `json:"translatedLanguage"`
	Uploader           string  `json:"uploader"`
	ExternalURL        *string `json:"externalUrl"`
	Version            int     `json:"version"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
	PublishAt          string  `json:"publishAt"`
}

// GetMangaChapters : Get a list of chapters for a manga.
// https://api.mangadex.org/docs.html#operation/get-manga-id-feed
func (s *ChapterService) GetMangaChapters(id string, params url.Values) (*ChapterList, error) {
	return s.GetMangaChaptersContext(context.Background(), id, params)
}

// GetMangaChaptersContext : GetMangaChapters with custom context.
func (s *ChapterService) GetMangaChaptersContext(ctx context.Context, id string, params url.Values) (*ChapterList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaChaptersPath, id)

	// Set request parameters
	u.RawQuery = params.Encode()

	var l ChapterList
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// ChapterReadMarkers : A response for getting a list of read chapters.
type ChapterReadMarkers struct {
	Result string   `json:"result"`
	Data   []string `json:"data"`
}

func (rmr *ChapterReadMarkers) GetResult() string {
	return rmr.Result
}

// GetReadMangaChapters : Get list of Chapter IDs that are marked as read for a specified manga ID.
// https://api.mangadex.org/docs.html#operation/get-manga-chapter-readmarkers
func (s *ChapterService) GetReadMangaChapters(id string) (*ChapterReadMarkers, error) {
	return s.GetReadMangaChaptersContext(context.Background(), id)
}

// GetReadMangaChaptersContext : GetReadMangaChapters with custom context.
func (s *ChapterService) GetReadMangaChaptersContext(ctx context.Context, id string) (*ChapterReadMarkers, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaReadMarkersPath, id)

	var rmr ChapterReadMarkers
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &rmr)
	return &rmr, err
}

// SetReadUnreadMangaChapters : Set read/unread manga chapters.
func (s *ChapterService) SetReadUnreadMangaChapters(id string, read, unRead []string) (*Response, error) {
	return s.SetReadUnreadMangaChaptersContext(context.Background(), id, read, unRead)
}

// SetReadUnreadMangaChaptersContext : SetReadUnreadMangaChapters with custom context.
func (s *ChapterService) SetReadUnreadMangaChaptersContext(ctx context.Context, id string, read, unRead []string) (*Response, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = fmt.Sprintf(MangaReadMarkersPath, id)

	// Set request body.
	req := map[string][]string{
		"chapterIdsRead":   read,
		"chapterIdsUnread": unRead,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var r Response
	err = s.client.RequestAndDecode(ctx, http.MethodPost, u.String(), bytes.NewBuffer(rBytes), &r)
	return &r, err
}
