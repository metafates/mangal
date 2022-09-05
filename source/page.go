package source

import (
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/util"
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
	// Size of the page in bytes
	Size uint64 `json:"-"`
	// Contents of the page
	Contents io.ReadCloser `json:"-"`
	// Chapter that the page belongs to.
	Chapter *Chapter `json:"-"`
}

// Download Page contents.
func (p *Page) Download() error {
	if p.URL == "" {
		log.Warnf("Page #%d has no URL", p.Index)
		return nil
	}

	log.Debugf("Downloading page #%d (%s)", p.Index, p.URL)

	req, err := http.NewRequest(http.MethodGet, p.URL, nil)
	if err != nil {
		log.Error(err)
		return err
	}

	req.Header.Set("Referer", p.Chapter.URL)
	req.Header.Set("User-Agent", constant.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New("http error: " + resp.Status)
		log.Error(err)
		return err
	}

	if resp.ContentLength == 0 {
		err = errors.New("http error: nothing was returned")
		log.Error(err)
		return err
	}

	p.Contents = resp.Body
	p.Size = uint64(util.Max(resp.ContentLength, 0))

	log.Debugf("Page #%d downloaded - %s", p.Index, humanize.Bytes(p.Size))
	return nil
}

// Close closes the page contents.
func (p *Page) Close() error {
	return p.Contents.Close()
}

// Read reads from the page contents.
func (p *Page) Read(b []byte) (int, error) {
	log.Debugf("Reading page contents #%d", p.Index)
	if p.Contents == nil {
		return 0, errors.New("page not downloaded")
	}

	return p.Contents.Read(b)
}

// Filename generates a filename for the page.
func (p *Page) Filename() (filename string) {
	filename = fmt.Sprintf("%d%s", p.Index, p.Extension)
	filename = util.PadZero(filename, 10)

	return
}

func (p *Page) Source() Source {
	return p.Chapter.Source()
}
