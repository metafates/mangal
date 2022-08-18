package source

import (
	"errors"
	"github.com/metafates/mangal/constant"
	"github.com/samber/lo"
	"io"
	"net/http"
)

type Page struct {
	URL       string `json:"url"`
	Index     uint16 `json:"index"`
	Extension string `json:"extension"`
	SourceID  string `json:"source_id"`
	Size      uint64
	Contents  io.ReadCloser
	Chapter   *Chapter
}

func (p *Page) Download() error {
	if p.URL == "" {
		return nil
	}

	req, err := http.NewRequest("GET", p.URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Referer", p.Chapter.URL)
	req.Header.Set("User-Agent", constant.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("http error: " + resp.Status)
	}

	if resp.ContentLength == 0 {
		return errors.New("http error: nothing was returned")
	}

	p.Contents = resp.Body
	p.Size = lo.Max([]uint64{uint64(resp.ContentLength), 0})

	return nil
}

func (p *Page) Close() error {
	return p.Contents.Close()
}

func (p *Page) Read(b []byte) (int, error) {
	if p.Contents == nil {
		return 0, errors.New("page not downloaded")
	}

	return p.Contents.Read(b)
}
