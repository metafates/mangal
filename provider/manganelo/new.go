package manganelo

import (
	"github.com/gocolly/colly"
	"github.com/metafates/mangal/constants"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"time"
)

func New() source.Source {
	manganelo := Manganelo{
		mangas:   make(map[string][]*source.Manga),
		chapters: make(map[string][]*source.Chapter),
		pages:    make(map[string][]*source.Page),
	}

	cacheDir := filepath.Join(lo.Must(os.UserCacheDir()), constants.CachePrefix)

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
		r.Headers.Set("User-Agent", constants.UserAgent)
	})

	// Get mangas
	mangasCollector.OnHTML(".search-story-item a.item-title", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		manga := source.Manga{
			Name:     e.Text,
			URL:      e.Request.AbsoluteURL(link),
			Index:    uint16(e.Index),
			Chapters: make([]*source.Chapter, 0),
			SourceID: manganelo.ID(),
		}

		manganelo.mangas[path] = append(manganelo.mangas[path], &manga)
	})

	_ = mangasCollector.Limit(&colly.LimitRule{
		Parallelism: 50,
		RandomDelay: 300 * time.Millisecond,
		DomainGlob:  "*",
	})

	chaptersCollector := baseCollector.Clone()
	chaptersCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", r.Ctx.GetAny("manga").(*source.Manga).URL)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("Host", "https://ww5.manganelo.tv/")
		r.Headers.Set("User-Agent", constants.UserAgent)
	})

	// Get chapters
	chaptersCollector.OnHTML(".chapter-name", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		manga := e.Request.Ctx.GetAny("manga").(*source.Manga)
		chapter := source.Chapter{
			Name:     e.Text,
			URL:      e.Request.AbsoluteURL(link),
			Index:    uint16(e.Index),
			Pages:    make([]*source.Page, 0),
			Manga:    manga,
			SourceID: manganelo.ID(),
		}
		manga.Chapters = append(manga.Chapters, &chapter)

		manganelo.chapters[path] = append(manganelo.chapters[path], &chapter)
	})
	_ = chaptersCollector.Limit(&colly.LimitRule{
		Parallelism: 50,
		RandomDelay: 300 * time.Millisecond,
		DomainGlob:  "*",
	})

	pagesCollector := baseCollector.Clone()
	pagesCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", r.Ctx.GetAny("chapter").(*source.Chapter).URL)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("User-Agent", constants.UserAgent)
	})

	// Get pages
	pagesCollector.OnHTML(".container-chapter-reader img", func(e *colly.HTMLElement) {
		link := e.Attr("data-src")
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		chapter := e.Request.Ctx.GetAny("chapter").(*source.Chapter)
		page := source.Page{
			URL:       link,
			Index:     uint16(e.Index),
			Chapter:   chapter,
			Extension: ".jpg",
			SourceID:  manganelo.ID(),
		}
		chapter.Pages = append(chapter.Pages, &page)

		manganelo.pages[path] = append(manganelo.pages[path], &page)
	})
	_ = pagesCollector.Limit(&colly.LimitRule{
		Parallelism: 50,
		RandomDelay: 300 * time.Millisecond,
		DomainGlob:  "*",
	})

	manganelo.mangasCollector = mangasCollector
	manganelo.chaptersCollector = chaptersCollector
	manganelo.pagesCollector = pagesCollector

	return &manganelo
}
