package scraper

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Source struct {
	// Base url address of the scraper
	Base string
	// Search page url address
	Search string
	// MangaAnchor css selector
	MangaAnchor string `toml:"manga_anchor"`
	// MangaTitle css selector
	MangaTitle string `toml:"manga_title"`
	// ChapterAnchor css selector
	ChapterAnchor string `toml:"chapter_anchor"`
	// ChapterTitle css selector
	ChapterTitle string `toml:"chapter_title"`
	// ChapterPages css selector
	ReaderPages string `toml:"reader_pages"`

	cache map[string][]*URL
}

func (s Source) paired(where URL, anchorSelector string, titleSelector string) ([]*URL, error) {
	doc, err := where.document()

	if err != nil {
		return nil, err
	}

	var urls []*URL

	anchors, titles := doc.Find(anchorSelector), doc.Find(titleSelector)

	anchors.Each(func(i int, anchor *goquery.Selection) {
		// It could be that due to the wrong selector there are more anchors than titles
		if i >= titles.Length() {
			return
		}

		// We need to convert node to selection to use Selection.Text() (textContent)
		// I don't know if there is a better way to do it...
		title := goquery.Selection{Nodes: []*html.Node{titles.Get(i)}}

		href, hasHref := anchor.Attr("href")
		if !hasHref {
			return
		}

		if !strings.HasPrefix(href, s.Base) {
			href = s.Base + href
		}

		url := URL{Address: href, Info: strings.TrimSpace(title.Text()), Source: &s}
		urls = append(urls, &url)
	})

	return urls, nil
}

// Mangas searches for mangas with given title
func (s Source) Mangas(title string) ([]*URL, error) {
	if s.cache == nil {
		s.cache = make(map[string][]*URL)
	}

	if val, ok := s.cache[title]; ok {
		return val, nil
	}

	search := URL{Address: fmt.Sprintf(s.Search, title), Source: &s}
	res, err := s.paired(search, s.MangaAnchor, s.MangaTitle)

	if err != nil {
		return nil, err
	}

	s.cache[title] = res
	return res, nil
}

// Chapters gets chapters of the given manga
func (s Source) Chapters(manga URL) ([]*URL, error) {
	if s.cache == nil {
		s.cache = make(map[string][]*URL)
	}

	if val, ok := s.cache[manga.Address]; ok {
		return val, nil
	}

	res, err := s.paired(manga, s.ChapterAnchor, s.ChapterTitle)

	if err != nil {
		return nil, err
	}

	s.cache[manga.Address] = res
	return res, nil
}

// Pages gets chapter pages
func (s Source) Pages(chapter URL) ([]*URL, error) {
	if s.cache == nil {
		s.cache = make(map[string][]*URL)
	}

	if val, ok := s.cache[chapter.Address]; ok {
		return val, nil
	}

	doc, err := chapter.document()

	if err != nil {
		return nil, err
	}

	var urls []*URL
	pages := doc.Find(s.ReaderPages)

	pages.Each(func(i int, page *goquery.Selection) {
		dataSrc, hasDataSrc := page.Attr("data-src")

		if !hasDataSrc {
			src, hasSrc := page.Attr("src")
			if !hasSrc {
				return
			}
			dataSrc = src
		}

		name := filepath.Base(dataSrc)
		url := URL{Address: dataSrc, Info: name, Source: &s}
		urls = append(urls, &url)
	})

	s.cache[chapter.Address] = urls
	return urls, nil
}
