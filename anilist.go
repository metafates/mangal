package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	levenshtein "github.com/ka-weihe/fast-levenshtein"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type AnilistURL struct {
	Address string `json:"address"`
	Title   string `json:"title"`
	ID      int    `json:"id"`
}

type AnilistPreferences struct {
	Connections map[string]*AnilistURL `json:"connections"`
	Token       struct {
		JWT       string `json:"jwt"`
		ExpiresAt int64  `json:"expires_at"`
	} `json:"token"`
}

type AnilistClient struct {
	ID          string
	Secret      string
	Preferences *AnilistPreferences `json:"preferences"`
}

// NewAnilistClient creates a new client for anilist integration
func NewAnilistClient(id, secret string) (*AnilistClient, error) {
	client := &AnilistClient{
		Preferences: &AnilistPreferences{
			Connections: make(map[string]*AnilistURL),
		},
		ID:     id,
		Secret: secret,
	}

	err := client.LoadPreferences()

	if err != nil {
		return nil, err
	}

	return client, nil
}

// AuthURL returns the URL for the Anilist authorization page
func (a *AnilistClient) AuthURL() string {
	return "https://anilist.co/api/v2/oauth/authorize?client_id=" + a.ID + "&response_type=code&redirect_uri=https://anilist.co/api/v2/oauth/pin"
}

// IsExpired returns true if the token is expired
func (a *AnilistClient) IsExpired() bool {
	// shortcut when token is not set
	if a.Preferences.Token.ExpiresAt == 0 {
		return true
	}

	return time.Now().UnixNano() > a.Preferences.Token.ExpiresAt
}

// Login to Anilist
func (a *AnilistClient) Login(code string) error {
	// anilist body for login
	body := map[string]interface{}{
		"client_id":     a.ID,
		"client_secret": a.Secret,
		"grant_type":    "authorization_code",
		"redirect_uri":  "https://anilist.co/api/v2/oauth/pin",
		"code":          code,
	}

	// encode body
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// create request
	req, err := http.NewRequest("POST", "https://anilist.co/api/v2/oauth/token", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	// set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// send request
	resp, err := http.DefaultClient.Do(req)

	// check for error
	if err != nil {
		return err
	}

	// check response code
	if resp.StatusCode != 200 {
		fmt.Println(resp)
		return errors.New("invalid response code " + strconv.Itoa(resp.StatusCode))
	}

	// decode response
	var response struct {
		AccessToken string `json:"access_token"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	// set token that would expire in a day
	a.Preferences.Token.JWT = response.AccessToken
	a.Preferences.Token.ExpiresAt = time.Now().UnixNano() + int64(time.Hour*24)

	// save preferences
	if err = a.SavePreferences(); err != nil {
		return err
	}

	return nil
}

// SavePreferences saves the preferences to the file
func (a *AnilistClient) SavePreferences() error {
	// get anilist file
	anilistFile, err := AnilistCacheFile()

	if err != nil {
		return err
	}

	// check if parent directory exists and create if not
	if exists, err := Afero.Exists(filepath.Dir(anilistFile)); err != nil {
		return err
	} else if !exists {
		if err := Afero.MkdirAll(filepath.Dir(anilistFile), 0777); err != nil {
			return err
		}
	}

	// check if file exists
	if exists, err := Afero.Exists(anilistFile); err != nil {
		return err
	} else if !exists {
		// create file and write current preferences
		jsonPreferences, err := json.Marshal(a.Preferences)
		if err != nil {
			return err
		}

		if err := Afero.WriteFile(anilistFile, jsonPreferences, 0777); err != nil {
			return err
		}

		return nil
	}

	// read file
	file, err := Afero.ReadFile(anilistFile)
	if err != nil {
		return err
	}

	// decode file that was read
	var preferences AnilistPreferences
	if err := json.NewDecoder(bytes.NewReader(file)).Decode(&preferences); err != nil {
		return err
	}

	// update preferences
	preferences.Connections = a.Preferences.Connections
	preferences.Token = a.Preferences.Token

	// write preferences to file
	jsonPreferences, err := json.Marshal(preferences)
	if err != nil {
		return err
	}

	if err := Afero.WriteFile(anilistFile, jsonPreferences, 0777); err != nil {
		return err
	}

	return nil
}

// LoadPreferences loads the preferences from the file
func (a *AnilistClient) LoadPreferences() error {
	// get anilist file
	anilistFile, err := AnilistCacheFile()

	if err != nil {
		return err
	}

	// check if file exists
	if exists, err := Afero.Exists(anilistFile); err != nil {
		return err
	} else if !exists {
		return nil
	}

	// read file
	file, err := Afero.ReadFile(anilistFile)
	if err != nil {
		return err
	}

	// decode file
	if err := json.NewDecoder(bytes.NewReader(file)).Decode(a.Preferences); err != nil {
		return err
	}

	return nil
}

// SearchManga searches for a manga
func (a *AnilistClient) SearchManga(manga string) ([]*AnilistURL, error) {
	// query to search for manga
	query := `
		query ($query: String) {
			Page (page: 1, perPage: 30) {
				media (search: $query, type: MANGA) {
					id
					title {
						romaji
						english
						native
					}
				}
			}
		}
`

	// prepare body
	body := map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"query": manga,
		},
	}

	// parse body to json
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// send request
	resp, err := http.Post(
		"https://graphql.anilist.co",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return nil, err
	}

	// decode response
	var response struct {
		Data struct {
			Page struct {
				Media []struct {
					ID    int `json:"id"`
					Title struct {
						Romaji  string `json:"romaji"`
						English string `json:"english"`
						Native  string `json:"native"`
					} `json:"title"`
				} `json:"media"`
			} `json:"page"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// convert to AnilistURL
	var urls = make([]*AnilistURL, len(response.Data.Page.Media))
	for i, media := range response.Data.Page.Media {
		urls[i] = &AnilistURL{
			ID:      media.ID,
			Address: "https://anilist.co/manga/" + ToString(media.ID),
			Title:   IfElse(media.Title.English != "", media.Title.English, media.Title.Romaji),
		}
	}

	return urls, nil
}

// ToAnilistURL will find manga on anilist similar to the given title
// It will cache the results
func (a *AnilistClient) ToAnilistURL(manga *URL) *AnilistURL {
	if cached, ok := a.Preferences.Connections[manga.Address]; ok {
		return cached
	}

	// search for manga
	urls, err := a.SearchManga(manga.Info)
	if err != nil {
		return nil
	}

	// find most similar manga using levenshtein distance
	var (
		closest         *AnilistURL
		closestDistance = 999
	)

	for _, url := range urls {
		distance := levenshtein.Distance(manga.Info, url.Title)
		if distance < closestDistance {
			closest = url
			closestDistance = distance
		}
	}

	// cache result
	a.Preferences.Connections[manga.Address] = closest
	_ = a.SavePreferences()

	return closest
}

// MarkChapter marks a chapter as read
func (a *AnilistClient) MarkChapter(manga *URL, chapter int) error {

	anilistManga := a.ToAnilistURL(manga)
	if anilistManga == nil {
		return errors.New("manga not found")
	}

	// query to mark chapter
	query := `
		mutation ($id: Int, $progress: Int) {
			SaveMediaListEntry (mediaId: $id, progress: $progress, status: CURRENT) {
				id
			}
		}
`

	// prepare body
	body := map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"id":       anilistManga.ID,
			"progress": chapter,
		},
	}

	// parse body to json
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// make request
	req, err := http.NewRequest(
		"POST",
		"https://graphql.anilist.co",
		bytes.NewBuffer(jsonBody),
	)

	// set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.Preferences.Token.JWT)
	req.Header.Set("Accept", "application/json")

	// send request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	// decode response
	var response struct {
		Data struct {
			SaveMediaListEntry struct {
				ID int `json:"id"`
			} `json:"SaveMediaListEntry"`
		} `json:"data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	return nil
}
