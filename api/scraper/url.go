package scraper

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"time"
)

type URL struct {
	// Address of this URL
	Address string
	// Info about this URL (e.g. page title)
	Info string
	// Source parent of this URL
	Source *Source
}

// Get http response from url
func (u URL) Get() (*http.Response, error) {
	client := http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("GET", u.Address, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", getRandomUserAgent())

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Cannot get '%s'. Status %d - %s", u.Address, resp.StatusCode, resp.Status))
	}

	return resp, nil
}

// document from URL
func (u URL) document() (*goquery.Document, error) {
	resp, err := u.Get()

	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal("Unexpected error while closing http body")
		}
	}()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, err
	}

	return doc, nil
}
