package source

import (
	"errors"
	"github.com/metafates/mangal/constant"
	"github.com/samber/lo"
	"io"
	"net/http"
)

// Page represents a page in a chapter
type Page struct {
	// URL of the page. Used to download the page.
	URL string
	// Index of the page in the chapter.
	Index uint16
	// Extension of the page image.
	Extension string
	// SourceID of the source the page is from.
	SourceID string
	// Size of the page in bytes
	Size uint64
	// Contents of the page
	Contents io.ReadCloser
	// Chapter that the page belongs to.
	Chapter *Chapter
}

// Download Page contents.
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

// Close closes the page contents.
func (p *Page) Close() error {
	return p.Contents.Close()
}

// Read reads from the page contents.
func (p *Page) Read(b []byte) (int, error) {
	if p.Contents == nil {
		return 0, errors.New("page not downloaded")
	}

	return p.Contents.Read(b)
}
