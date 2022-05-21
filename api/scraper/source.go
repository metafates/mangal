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
	// ChapterPanels css selector
	ChapterPanels string `toml:"chapter_panels"`
}

func (s Source) paired(where URL, anchorSelector string, titleSelector string) ([]*URL, error) {
	doc, err := where.document()

	if err != nil {
		return nil, err
	}

	var urls []*URL

	anchors, titles := doc.Find(anchorSelector), doc.Find(titleSelector)

	anchors.Each(func(i int, anchor *goquery.Selection) {
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
	search := URL{Address: fmt.Sprintf(s.Search, title), Source: &s}
	return s.paired(search, s.MangaAnchor, s.MangaTitle)
}

// Chapters gets chapters of the given manga
func (s Source) Chapters(manga URL) ([]*URL, error) {
	return s.paired(manga, s.ChapterAnchor, s.ChapterTitle)
}

// Panels gets chapter panels
func (s Source) Panels(chapter URL) ([]*URL, error) {
	doc, err := chapter.document()

	if err != nil {
		return nil, err
	}

	var urls []*URL
	panels := doc.Find(s.ChapterPanels)

	panels.Each(func(i int, panel *goquery.Selection) {
		dataSrc, hasDataSrc := panel.Attr("data-src")

		if !hasDataSrc {
			src, hasSrc := panel.Attr("src")
			if !hasSrc {
				return
			}
			dataSrc = src
		}

		name := filepath.Base(dataSrc)
		url := URL{Address: dataSrc, Info: name, Source: &s}
		urls = append(urls, &url)
	})

	return urls, nil
}
