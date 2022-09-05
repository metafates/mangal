package manganelo

import (
	"github.com/gocolly/colly"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	delay       = time.Millisecond * 500
	parallelism = 50

	mangasSelector   = ".search-story-item a.item-title"
	chaptersSelector = ".chapter-name"
	pageSelector     = ".container-chapter-reader img"
)

func New() source.Source {
	manganelo := Manganelo{
		mangas:   make(map[string][]*source.Manga),
		chapters: make(map[string][]*source.Chapter),
		pages:    make(map[string][]*source.Page),
	}

	cacheDir := filepath.Join(lo.Must(os.UserCacheDir()), constant.CachePrefix)

	collectorOptions := []func(*colly.Collector){
		colly.AllowURLRevisit(),
		colly.Async(true),
		colly.CacheDir(cacheDir),
	}

	baseCollector := colly.NewCollector(collectorOptions...)
	baseCollector.SetRequestTimeout(20 * time.Second)

	mangasCollector := baseCollector.Clone()
	mangasCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", "https://google.com")
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("Host", "https://ww5.manganelo.tv/")
		r.Headers.Set("User-Agent", constant.UserAgent)
	})

	// Get mangas
	mangasCollector.OnHTML(mangasSelector, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		path := e.Request.URL.String()
		url := e.Request.AbsoluteURL(link)
		manga := source.Manga{
			Name:     e.Text,
			URL:      url,
			Index:    uint16(e.Index),
			Chapters: make([]*source.Chapter, 0),
			ID:       filepath.Base(url),
			Source:   &manganelo,
		}

		manganelo.mangas[path] = append(manganelo.mangas[path], &manga)
	})

	_ = mangasCollector.Limit(&colly.LimitRule{
		Parallelism: parallelism,
		RandomDelay: delay,
		DomainGlob:  "*",
	})

	chaptersCollector := baseCollector.Clone()
	chaptersCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", r.Ctx.GetAny("manga").(*source.Manga).URL)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("Host", "https://ww5.manganelo.tv/")
		r.Headers.Set("User-Agent", constant.UserAgent)
	})

	// Get chapters
	chaptersCollector.OnHTML(chaptersSelector, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		manga := e.Request.Ctx.GetAny("manga").(*source.Manga)
		url := e.Request.AbsoluteURL(link)

		var (
			volume string
			name   = e.Text
		)

		if strings.HasPrefix(name, "Vol.") {
			splitted := strings.Split(name, " ")
			volume = splitted[0]
			name = strings.Join(splitted[1:], " ")
		}

		chapter := source.Chapter{
			Name:   name,
			URL:    url,
			Index:  uint16(e.Index),
			Pages:  make([]*source.Page, 0),
			ID:     filepath.Base(url),
			Manga:  manga,
			Volume: volume,
		}
		manga.Chapters = append(manga.Chapters, &chapter)

		manganelo.chapters[path] = append(manganelo.chapters[path], &chapter)
	})
	_ = chaptersCollector.Limit(&colly.LimitRule{
		Parallelism: parallelism,
		RandomDelay: delay,
		DomainGlob:  "*",
	})

	pagesCollector := baseCollector.Clone()
	pagesCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", r.Ctx.GetAny("chapter").(*source.Chapter).URL)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("User-Agent", constant.UserAgent)
	})

	// Get pages
	pagesCollector.OnHTML(pageSelector, func(e *colly.HTMLElement) {
		link := e.Attr("data-src")
		ext := filepath.Ext(link)
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		chapter := e.Request.Ctx.GetAny("chapter").(*source.Chapter)
		page := source.Page{
			URL:       link,
			Index:     uint16(e.Index),
			Chapter:   chapter,
			Extension: ext,
		}
		chapter.Pages = append(chapter.Pages, &page)

		manganelo.pages[path] = append(manganelo.pages[path], &page)
	})
	_ = pagesCollector.Limit(&colly.LimitRule{
		Parallelism: parallelism,
		RandomDelay: delay,
		DomainGlob:  "*",
	})

	manganelo.mangasCollector = mangasCollector
	manganelo.chaptersCollector = chaptersCollector
	manganelo.pagesCollector = pagesCollector

	return &manganelo
}
