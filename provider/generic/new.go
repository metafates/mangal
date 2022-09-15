package generic

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/where"
	"path/filepath"
	"time"
)

func New(conf *Configuration) source.Source {
	s := Scraper{
		mangas:   make(map[string][]*source.Manga),
		chapters: make(map[string][]*source.Chapter),
		pages:    make(map[string][]*source.Page),
		config:   conf,
	}

	collectorOptions := []func(*colly.Collector){
		colly.AllowURLRevisit(),
		colly.Async(true),
		colly.CacheDir(where.Cache()),
	}

	baseCollector := colly.NewCollector(collectorOptions...)
	baseCollector.SetRequestTimeout(20 * time.Second)

	mangasCollector := baseCollector.Clone()
	mangasCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", "https://google.com")
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("Host", s.config.BaseURL)
		r.Headers.Set("User-Agent", constant.UserAgent)
	})

	// Get mangas
	mangasCollector.OnHTML("html", func(e *colly.HTMLElement) {
		elements := e.DOM.Find(s.config.MangaExtractor.Selector)
		path := e.Request.URL.String()
		s.mangas[path] = make([]*source.Manga, elements.Length())

		elements.Each(func(i int, selection *goquery.Selection) {
			link := s.config.MangaExtractor.URL(selection)
			url := e.Request.AbsoluteURL(link)
			manga := source.Manga{
				Name:     s.config.MangaExtractor.Name(selection),
				URL:      url,
				Index:    uint16(e.Index),
				Chapters: make([]*source.Chapter, 0),
				ID:       filepath.Base(url),
				Source:   &s,
			}
			manga.Metadata.Cover = s.config.MangaExtractor.Cover(selection)

			s.mangas[path][i] = &manga
		})
	})

	_ = mangasCollector.Limit(&colly.LimitRule{
		Parallelism: int(s.config.Parallelism),
		RandomDelay: s.config.Delay,
		DomainGlob:  "*",
	})

	chaptersCollector := baseCollector.Clone()
	chaptersCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", r.Ctx.GetAny("manga").(*source.Manga).URL)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("Host", s.config.BaseURL)
		r.Headers.Set("User-Agent", constant.UserAgent)
	})

	// Get chapters
	chaptersCollector.OnHTML("html", func(e *colly.HTMLElement) {
		elements := e.DOM.Find(s.config.ChapterExtractor.Selector)
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		s.chapters[path] = make([]*source.Chapter, elements.Length())
		manga := e.Request.Ctx.GetAny("manga").(*source.Manga)
		manga.Chapters = make([]*source.Chapter, elements.Length())

		elements.Each(func(i int, selection *goquery.Selection) {
			link := s.config.ChapterExtractor.URL(selection)
			url := e.Request.AbsoluteURL(link)

			chapter := source.Chapter{
				Name:   s.config.ChapterExtractor.Name(selection),
				URL:    url,
				Index:  uint16(e.Index),
				Pages:  make([]*source.Page, 0),
				ID:     filepath.Base(url),
				Manga:  manga,
				Volume: s.config.ChapterExtractor.Volume(selection),
			}
			manga.Chapters[i] = &chapter
			s.chapters[path][i] = &chapter
		})
	})
	_ = chaptersCollector.Limit(&colly.LimitRule{
		Parallelism: int(s.config.Parallelism),
		RandomDelay: s.config.Delay,
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
		elements := e.DOM.Find(s.config.PageExtractor.Selector)
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		s.pages[path] = make([]*source.Page, elements.Length())
		chapter := e.Request.Ctx.GetAny("chapter").(*source.Chapter)
		chapter.Pages = make([]*source.Page, elements.Length())

		elements.Each(func(i int, selection *goquery.Selection) {
			link := s.config.PageExtractor.URL(selection)
			ext := filepath.Ext(link)
			page := source.Page{
				URL:       link,
				Index:     uint16(i),
				Chapter:   chapter,
				Extension: ext,
			}
			chapter.Pages[i] = &page
			s.pages[path][i] = &page
		})
	})
	_ = pagesCollector.Limit(&colly.LimitRule{
		Parallelism: int(s.config.Parallelism),
		RandomDelay: s.config.Delay,
		DomainGlob:  "*",
	})

	s.mangasCollector = mangasCollector
	s.chaptersCollector = chaptersCollector
	s.pagesCollector = pagesCollector

	return &s
}
