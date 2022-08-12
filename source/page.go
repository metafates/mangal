package source

import (
	"errors"
	"github.com/metafates/mangal/constant"
	lua "github.com/yuin/gopher-lua"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Page struct {
	URL       string `json:"url"`
	Index     uint16 `json:"index"`
	Extension string `json:"extension"`
	SourceID  string `json:"source_id"`
	Contents  io.ReadCloser
	Chapter   *Chapter
}

func pageFromTable(table *lua.LTable, chapter *Chapter) (*Page, error) {
	url := table.RawGetString("url")

	if url.Type() != lua.LTString {
		return nil, errors.New("type of field \"url\" should be string")
	}

	index := table.RawGetString("index")

	if index.Type() != lua.LTNumber {
		return nil, errors.New("type of field \"index\" should be number")
	}

	num, err := strconv.ParseUint(index.String(), 10, 16)

	if err != nil {
		return nil, errors.New("index must be an unsigned 16 bit integer")
	}

	page := &Page{
		URL:       strings.TrimSpace(url.String()),
		Index:     uint16(num),
		Chapter:   chapter,
		Extension: ".jpg",
	}

	chapter.Pages = append(chapter.Pages, page)
	return page, nil
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
