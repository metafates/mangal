package manganelo

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/source"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"regexp"
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
	mangasCollector.OnHTML("html", func(e *colly.HTMLElement) {
		elements := e.DOM.Find(mangasSelector)
		path := e.Request.URL.String()
		manganelo.mangas[path] = make([]*source.Manga, elements.Length())

		elements.Each(func(i int, selection *goquery.Selection) {
			link, _ := selection.Attr("href")
			url := e.Request.AbsoluteURL(link)
			manga := source.Manga{
				Name:     selection.Text(),
				URL:      url,
				Index:    uint16(e.Index),
				Chapters: make([]*source.Chapter, 0),
				ID:       filepath.Base(url),
				Source:   &manganelo,
			}

			manganelo.mangas[path][i] = &manga
		})
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
	chaptersCollector.OnHTML("html", func(e *colly.HTMLElement) {
		elements := e.DOM.Find(chaptersSelector)
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		manganelo.chapters[path] = make([]*source.Chapter, elements.Length())
		manga := e.Request.Ctx.GetAny("manga").(*source.Manga)
		manga.Chapters = make([]*source.Chapter, elements.Length())

		summary, _ := e.DOM.Find("#panel-story-info-description").Html()
		// remove html tags
		manga.Summary = regexp.MustCompile(`<.*?>`).ReplaceAllString(summary, "")

		genres := e.
			DOM.
			Find("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-info > div.story-info-right > table > tbody > tr:nth-child(4) > td.table-value a").
			Map(func(_ int, selection *goquery.Selection) string {
				return selection.Text()
			})
		manga.Genres = genres

		author := e.
			DOM.
			Find("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-info > div.story-info-right > table > tbody > tr:nth-child(2) > td.table-value > a").
			Text()
		manga.Author = author

		elements.Each(func(i int, selection *goquery.Selection) {
			link, _ := selection.Attr("href")
			url := e.Request.AbsoluteURL(link)

			var (
				volume string
				name   = selection.Text()
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
			manga.Chapters[i] = &chapter
			manganelo.chapters[path][i] = &chapter
		})
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
	pagesCollector.OnHTML("html", func(e *colly.HTMLElement) {
		elements := e.DOM.Find(pageSelector)
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		manganelo.pages[path] = make([]*source.Page, elements.Length())
		chapter := e.Request.Ctx.GetAny("chapter").(*source.Chapter)
		chapter.Pages = make([]*source.Page, elements.Length())

		elements.Each(func(i int, selection *goquery.Selection) {
			link, _ := selection.Attr("data-src")
			ext := filepath.Ext(link)
			page := source.Page{
				URL:       link,
				Index:     uint16(i),
				Chapter:   chapter,
				Extension: ext,
			}
			chapter.Pages[i] = &page
			manganelo.pages[path][i] = &page
		})

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
