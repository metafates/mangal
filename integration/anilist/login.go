package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/network"
	"github.com/samber/lo"
	"net/http"
	"strconv"
)

// login to Anilist
func (a *Anilist) login() error {
	log.Info("Logging in to Anilist")

	if a.id() == "" {
		e := fmt.Errorf("no ID set")
		log.Error(e)
		return e
	}
	if a.secret() == "" {
		e := fmt.Errorf("no secret set")
		log.Error(e)
		return e
	}
	if a.code() == "" {
		e := fmt.Errorf("no code set")
		log.Error(e)
		return e
	}

	// anilist body for login
	body := map[string]interface{}{
		"client_id":     a.id(),
		"client_secret": a.secret(),
		"grant_type":    "authorization_code",
		"redirect_uri":  "https://anilist.co/api/v2/oauth/pin",
		"code":          a.code(),
	}

	// encode body
	jsonBody := lo.Must(json.Marshal(&body))

	// create request
	log.Info("Sending login request to Anilist")
	req, err := http.NewRequest(http.MethodPost, "https://anilist.co/api/v2/oauth/token", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error(err)
		return err
	}

	// set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// send request
	resp, err := network.Client.Do(req)

	// check for error
	if err != nil {
		log.Error(err)
		return err
	}

	// check response code
	if resp.StatusCode != http.StatusOK {
		log.Info("Request failed with status code: " + strconv.Itoa(resp.StatusCode))
		return fmt.Errorf("invalid response code %d", resp.StatusCode)
	}

	// decode response
	log.Info("Decoding response from Anilist")
	var response struct {
		AccessToken string `json:"access_token"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error(err)
		return err
	}

	// set token
	log.Info("Logged in Anilist")
	a.token = response.AccessToken

	return nil
}
