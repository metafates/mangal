package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Scraper struct {
	Source            *Source
	MangaCollector    *colly.Collector
	ChaptersCollector *colly.Collector
	PagesCollector    *colly.Collector
	FilesCollector    *colly.Collector

	// Manga maps search url with manga urls
	Manga map[string][]*URL
	// Chapters maps manga url with chapters urls
	Chapters map[string][]*URL
	// Pages maps chapter url with pages urls
	Pages map[string][]*URL
	Files *RwMap[string, *bytes.Buffer]
}

// URL represents an url with relation to another url with useful information
type URL struct {
	Relation *URL     `json:"-"`
	Scraper  *Scraper `json:"-"`
	Address  string   `json:"address"`
	Info     string   `json:"info"`
	Index    int      `json:"index,omitempty"`
}

// MakeSourceScraper makes a scraper for a source
func MakeSourceScraper(source *Source) *Scraper {

	scraper := Scraper{
		Source: source,

		Manga:    make(map[string][]*URL),
		Chapters: make(map[string][]*URL),
		Pages:    make(map[string][]*URL),
		Files:    &RwMap[string, *bytes.Buffer]{data: make(map[string]*bytes.Buffer)},
	}

	var (
		collectorOptions []func(*colly.Collector)
		collector        *colly.Collector
	)

	collectorOptions = append(collectorOptions, colly.Async(true))
	defaultCacheDir, err := os.UserCacheDir()
	if err == nil {
		collectorOptions = append(collectorOptions, colly.CacheDir(filepath.Join(defaultCacheDir, CachePrefix)))
	}

	collector = colly.NewCollector(collectorOptions...)
	collector.SetRequestTimeout(20 * time.Second)
	collector.AllowURLRevisit = true

	extensions.RandomUserAgent(collector)

	// Manga collector
	mangaCollector := collector.Clone()

	// Prevent scraper from being blocked
	mangaCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", Referer)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("Host", source.Base)
	})

	// Get all manga urls
	mangaCollector.OnHTML(source.MangaAnchor, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		scraper.Manga[path] = append(scraper.Manga[path], &URL{Address: e.Request.AbsoluteURL(link), Scraper: &scraper})
	})

	// Get all manga titles
	mangaCollector.OnHTML(source.MangaTitle, func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.Text)
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		if e.Index < len(scraper.Manga[path]) {
			scraper.Manga[path][e.Index].Info = title
		}
	})

	_ = mangaCollector.Limit(&colly.LimitRule{
		Parallelism: Parallelism,
		RandomDelay: time.Duration(source.RandomDelayMs) * time.Millisecond,
		DomainGlob:  "*",
	})

	// Paths collector
	chaptersCollector := collector.Clone()
	chaptersCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", Referer)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("Host", source.Base)
	})
	chaptersCollector.OnHTML("html", func(html *colly.HTMLElement) {
		var urls []*URL

		// Get all chapters urls
		html.ForEach(source.ChapterAnchor, func(_ int, e *colly.HTMLElement) {
			link := e.Attr("href")
			path := e.Request.AbsoluteURL(e.Request.URL.Path)
			u := &URL{Address: e.Request.AbsoluteURL(link), Scraper: &scraper, Index: e.Index}
			urls = append(urls, u)
			scraper.Chapters[path] = append(scraper.Chapters[path], u)
		})

		// Get all chapter titles
		html.ForEachWithBreak(source.ChapterTitle, func(i int, e *colly.HTMLElement) bool {
			title := strings.TrimSpace(e.Text)
			path := e.Request.AbsoluteURL(e.Request.URL.Path)

			if e.Index >= len(scraper.Chapters[path]) {
				return false
			}

			scraper.Chapters[path][e.Index].Info = title
			return true
		})

		// Reverse chapters indexes
		if source.ChaptersReversed {
			length := len(urls)
			for _, u := range urls {
				u.Index = length - u.Index
			}
		}
	})
	_ = chaptersCollector.Limit(&colly.LimitRule{
		Parallelism: Parallelism,
		RandomDelay: time.Duration(source.RandomDelayMs) * time.Millisecond,
		DomainGlob:  "*",
	})

	// Pages collector
	pagesCollector := collector.Clone()
	pagesCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", Referer)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
	})
	pagesCollector.OnHTML(source.ReaderPage, func(e *colly.HTMLElement) {
		attributes := []string{
			"src",
			"data-src",
			"data-url",
			"href",
			"data-original",
			"data-original-src",
			"data-original-url",
			"data-srcset",
			"data-src-set",
		}

		attr, _ := Find(attributes, func(attr string) bool {
			return e.Attr(attr) != ""
		})

		link := e.Attr(attr)

		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		scraper.Pages[path] = append(scraper.Pages[path], &URL{Address: link, Scraper: &scraper, Index: e.Index})
	})
	_ = pagesCollector.Limit(&colly.LimitRule{
		Parallelism: Parallelism,
		RandomDelay: time.Duration(source.RandomDelayMs) * time.Millisecond,
		DomainGlob:  "*",
	})

	filesCollector := collector.Clone()
	filesCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", source.ChaptersBase)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	})
	filesCollector.OnResponse(func(r *colly.Response) {
		scraper.Files.Set(r.Request.AbsoluteURL(r.Request.URL.Path), bytes.NewBuffer(r.Body))
	})

	scraper.MangaCollector = mangaCollector
	scraper.ChaptersCollector = chaptersCollector
	scraper.PagesCollector = pagesCollector
	scraper.FilesCollector = filesCollector

	return &scraper
}

// SearchManga searches for manga by name
func (s *Scraper) SearchManga(title string) ([]*URL, error) {
	// lowercase titles will produce the same results but will be useful for caching
	query := strings.ReplaceAll(title, " ", s.Source.WhitespaceEscape)
	address := fmt.Sprintf(s.Source.SearchTemplate, url.QueryEscape(strings.TrimSpace(strings.ToLower(query))))

	if urls, ok := s.Manga[address]; ok {
		return urls, nil
	}

	err := s.MangaCollector.Visit(address)

	if err != nil {
		return nil, err
	}

	s.MangaCollector.Wait()
	return s.Manga[address], nil
}

// GetChapters returns manga chapters for given manga url
func (s *Scraper) GetChapters(manga *URL) ([]*URL, error) {
	if urls, ok := s.Chapters[manga.Address]; ok {
		return urls, nil
	}

	err := s.ChaptersCollector.Visit(manga.Address)

	if err != nil {
		return nil, err
	}

	s.ChaptersCollector.Wait()

	// Add relation to this manga url for each chapter
	// It shouldn't affect performance since there won't be more than 1000 chapters as worst case
	for _, chapter := range s.Chapters[manga.Address] {
		chapter.Relation = manga
	}

	return s.Chapters[manga.Address], nil
}

// GetPages returns manga pages for given chapter url
func (s *Scraper) GetPages(chapter *URL) ([]*URL, error) {
	if urls, ok := s.Pages[chapter.Address]; ok {
		return urls, nil
	}

	err := s.PagesCollector.Visit(chapter.Address)

	if err != nil {
		return nil, err
	}

	s.PagesCollector.Wait()

	// Add relation to this chapter url for each page
	// It shouldn't affect performance since there won't be more than ~400 pages as worst case (usually it's 30)
	for _, page := range s.Pages[chapter.Address] {
		page.Relation = chapter
	}

	return s.Pages[chapter.Address], nil
}

// GetFile returns manga file for given page url
func (s *Scraper) GetFile(file *URL) (*bytes.Buffer, error) {
	if data, ok := s.Files.Get(file.Address); ok {
		return data, nil
	}

	err := s.FilesCollector.Visit(file.Address)

	if err != nil {
		return nil, err
	}

	s.FilesCollector.Wait()

	data, ok := s.Files.Get(file.Address)

	if ok {
		return data, nil
	}

	return nil, errors.New("Couldn't get file at " + file.Address)
}

// ResetFiles resets scraper files cache
func (s *Scraper) ResetFiles() {
	s.Files.data = make(map[string]*bytes.Buffer)
}
