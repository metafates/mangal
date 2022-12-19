package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/network"
	"github.com/metafates/mangal/query"
	"github.com/samber/lo"
	"net/http"
	"strconv"
)

type searchByNameResponse struct {
	Data struct {
		Page struct {
			Media []*Manga `json:"media"`
		} `json:"page"`
	} `json:"data"`
}

type searchByIDResponse struct {
	Data struct {
		Media *Manga `json:"media"`
	} `json:"data"`
}

// GetByID returns the manga with the given id.
// If the manga is not found, it returns nil.
func GetByID(id int) (*Manga, error) {
	if manga := idCacher.Get(id); manga.IsPresent() {
		return manga.MustGet(), nil
	}

	// prepare body
	log.Infof("Searching anilist for manga with id: %d", id)
	body := map[string]interface{}{
		"query": searchByIDQuery,
		"variables": map[string]interface{}{
			"id": id,
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
	req, err := http.NewRequest(http.MethodPost, "https://graphql.anilist.co", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := network.Client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Anilist returned status code " + strconv.Itoa(resp.StatusCode))
		return nil, fmt.Errorf("invalid response code %d", resp.StatusCode)
	}

	// decode response
	var response searchByIDResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error(err)
		return nil, err
	}

	manga := response.Data.Media
	log.Infof("Got response from Anilist, found manga with id %d", manga.ID)
	_ = idCacher.Set(id, manga)
	return manga, nil
}

// SearchByName returns a list of mangas that match the given name.
func SearchByName(name string) ([]*Manga, error) {
	name = normalizedName(name)
	_ = query.Remember(name, 1)

	if _, failed := failCacher.Get(name).Get(); failed {
		return nil, fmt.Errorf("failed to search for %s", name)
	}

	if ids, ok := searchCacher.Get(name).Get(); ok {
		mangas := lo.FilterMap(ids, func(item, _ int) (*Manga, bool) {
			return idCacher.Get(item).Get()
		})

		if len(mangas) == 0 {
			_ = searchCacher.Delete(name)
			return SearchByName(name)
		}

		return mangas, nil
	}

	// prepare body
	log.Infof("Searching anilist for manga %s", name)
	body := map[string]any{
		"query": searchByNameQuery,
		"variables": map[string]any{
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
	req, err := http.NewRequest(http.MethodPost, "https://graphql.anilist.co", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := network.Client.Do(req)

	if err != nil {
		log.Error(err)
		_ = failCacher.Set(name, true)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Anilist returned status code " + strconv.Itoa(resp.StatusCode))
		_ = failCacher.Set(name, true)
		return nil, fmt.Errorf("invalid response code %d", resp.StatusCode)
	}

	// decode response
	var response searchByNameResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error(err)
		return nil, err
	}

	mangas := response.Data.Page.Media
	log.Infof("Got response from Anilist, found %d results", len(mangas))
	ids := make([]int, len(mangas))
	for i, manga := range mangas {
		ids[i] = manga.ID
		_ = idCacher.Set(manga.ID, manga)
	}
	_ = searchCacher.Set(name, ids)
	return mangas, nil
}
