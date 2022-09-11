package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/log"
	"net/http"
	"strconv"
)

var searchQuery = `
query ($query: String) {
	Page (page: 1, perPage: 30) {
		media (search: $query, type: MANGA) {
			id
			idMal
			title {
				romaji
				english
				native
			}
			description(asHtml: false)
			tags {
				name
			}
			genres
			coverImage {
				extraLarge
			}
			characters (page: 1, perPage: 10, role: MAIN) {
				nodes {
					name {
						full
					}
				}
			}
			startDate {
				year
				month	
				day
			}
			endDate {
				year
				month	
				day
			}
			status
			synonyms
			siteUrl
			countryOfOrigin
			externalLinks {
				url
			}
		}
	}
}
`

func Search(name string) ([]*Manga, error) {
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
				Media []*Manga `json:"media"`
			} `json:"page"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error(err)
		return nil, err
	}

	log.Info("Got response from Anilist, found " + strconv.Itoa(len(response.Data.Page.Media)) + " results")
	return response.Data.Page.Media, nil
}
