package mangodex

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

const (
	GetUserFollowedMangaListPath = "user/follows/manga"
	GetLoggedUserPath            = "user/me"
)

// UserService : Provides User services provided by the API.
type UserService service

// GetUserFollowedMangaList : Return list of followed Manga.
// https://api.mangadex.org/docs.html#operation/get-user-follows-manga
func (s *UserService) GetUserFollowedMangaList(limit, offset int, includes []string) (*MangaList, error) {
	return s.GetUserFollowedMangaListContext(context.Background(), limit, offset, includes)
}

// GetUserFollowedMangaListContext : GetUserFollowedMangaListPath with custom context.
func (s *UserService) GetUserFollowedMangaListContext(ctx context.Context, limit, offset int, includes []string) (*MangaList, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = GetUserFollowedMangaListPath

	// Set required query parameters
	q := u.Query()
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	for _, i := range includes {
		q.Add("includes[]", i)
	}
	u.RawQuery = q.Encode()

	var l MangaList
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &l)
	return &l, err
}

// UserResponse : Typical User response.
type UserResponse struct {
	Result   string `json:"result"`
	Response string `json:"response"`
	Data     User   `json:"data"`
}

func (ur *UserResponse) GetResult() string {
	return ur.Result
}

// User : Info on a MangaDex user.
type User struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Attributes    UserAttributes `json:"attributes"`
	Relationships []Relationship `json:"relationships"`
}

// UserAttributes : Attributes of a User.
type UserAttributes struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Version  int      `json:"version"`
}

// GetLoggedUser : Return logged UserResponse.
// https://api.mangadex.org/docs.html#operation/get-user-follows-group
func (s *UserService) GetLoggedUser() (*UserResponse, error) {
	return s.GetLoggedUserContext(context.Background())
}

// GetLoggedUserContext : GetLoggedUser with custom context.
func (s *UserService) GetLoggedUserContext(ctx context.Context) (*UserResponse, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = GetLoggedUserPath

	var r UserResponse
	err := s.client.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &r)
	return &r, err
}
