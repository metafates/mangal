package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/source"
	"net/http"
	"strconv"
)

var markReadQuery = `
mutation ($id: Int, $progress: Int) {
	SaveMediaListEntry (mediaId: $id, progress: $progress, status: CURRENT) {
		id
	}
}
`

func (a *Anilist) MarkRead(chapter *source.Chapter) error {
	if a.token == "" {
		err := a.login()
		if err != nil {
			log.Error(err)
			return err
		}
	}

	manga, err := a.findClosestMangaOnAnilist(chapter.Manga)
	if err != nil {
		log.Error(err)
		return err
	}

	// prepare body
	body := map[string]interface{}{
		"query": markReadQuery,
		"variables": map[string]interface{}{
			"id":       manga.id,
			"progress": chapter.Index,
		},
	}

	// parse body to json
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Error(err)
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
	req.Header.Set("Authorization", "Bearer "+a.token)
	req.Header.Set("Accept", "application/json")

	// send request
	log.Info("Sending request to Anilist: " + string(jsonBody))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Info("Request failed with status code: " + strconv.Itoa(resp.StatusCode))
		return fmt.Errorf("invalid response code %d", resp.StatusCode)
	}

	// decode response
	var response struct {
		Data struct {
			SaveMediaListEntry struct {
				ID int `json:"id"`
			} `json:"SaveMediaListEntry"`
		} `json:"data"`
	}

	return json.NewDecoder(resp.Body).Decode(&response)
}
