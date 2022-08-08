package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	levenshtein "github.com/ka-weihe/fast-levenshtein"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"net/http"
	"strconv"
)

type anilistManga struct {
	url  string
	name string
	id   int
}

var searchQuery = `
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

func (a *Anilist) searchManga(name string) ([]*anilistManga, error) {
	// prepare body
	log.Info("Searching anilist for manga: " + name)
	body := map[string]interface{}{
		"query": searchQuery,
		"variables": map[string]interface{}{
			"query": name,
		},
	}

	// parse body to json
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// send request
	log.Info("Sending request to Anilist")
	resp, err := http.Post(
		"https://graphql.anilist.co",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Anilist returned status code " + strconv.Itoa(resp.StatusCode))
		return nil, fmt.Errorf("invalid response code %d", resp.StatusCode)
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
		log.Error(err)
		return nil, err
	}

	// convert to AnilistURL
	log.Info("Got response from Anilist, found " + strconv.Itoa(len(response.Data.Page.Media)) + " results")
	var urls = make([]*anilistManga, len(response.Data.Page.Media))
	for i, media := range response.Data.Page.Media {
		name = media.Title.English
		if name == "" {
			name = media.Title.Romaji
		}

		log.Info("Found manga: " + name)
		urls[i] = &anilistManga{
			id:   media.ID,
			url:  "https://anilist.co/manga/" + strconv.Itoa(media.ID),
			name: name,
		}
	}

	return urls, nil
}

func (a *Anilist) findClosestMangaOnAnilist(manga *source.Manga) (*anilistManga, error) {
	if cached, ok := a.cache[manga.URL]; ok {
		log.Info("Found cached manga: " + cached.name)
		return cached, nil
	}

	// search for manga on anilist
	urls, err := a.searchManga(manga.Name)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if len(urls) == 0 {
		log.Warn("No results found on Anilist for manga" + manga.Name)
		return nil, nil
	}

	// find the closest match
	closest := lo.MinBy(urls, func(a, b *anilistManga) bool {
		return levenshtein.Distance(manga.Name, a.name) < levenshtein.Distance(manga.Name, b.name)
	})

	log.Info("Found closest match: " + closest.name)
	a.cache[manga.URL] = closest
	return closest, nil
}
