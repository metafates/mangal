package scraper

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/util"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ContextCollector[T any] struct {
	Collector *colly.Collector
	Ctx       T
}

type Scraper struct {
	Source            *Source
	mangaCollector    *colly.Collector
	chaptersCollector *ContextCollector[*URL]
	pagesCollector    *ContextCollector[*URL]
	filesCollector    *colly.Collector

	// Manga maps search url with manga urls
	manga map[string][]*URL
	// Chapters maps manga url with chapters urls
	chapters map[string][]*URL
	// Pages maps chapter url with pages urls
	pages map[string][]*URL

	Files *util.RwMap[string, *bytes.Buffer]
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

		manga:    make(map[string][]*URL),
		chapters: make(map[string][]*URL),
		pages:    make(map[string][]*URL),
		Files:    util.NewRwMap[string, *bytes.Buffer](),
	}

	var (
		collectorOptions []func(*colly.Collector)
		collector        *colly.Collector
	)

	collectorOptions = append(collectorOptions, colly.Async(true))
	defaultCacheDir, err := os.UserCacheDir()
	if err == nil {
		collectorOptions = append(collectorOptions, colly.CacheDir(filepath.Join(defaultCacheDir, common.CachePrefix)))
	}

	collector = colly.NewCollector(collectorOptions...)
	collector.AllowURLRevisit = true

	extensions.RandomUserAgent(collector)

	// Manga collector
	mangaCollector := collector.Clone()

	// Prevent scraper from being blocked
	mangaCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", common.Referer)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("Host", source.Base)
	})

	// Get all manga urls
	mangaCollector.OnHTML(source.MangaAnchor, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		scraper.manga[path] = append(scraper.manga[path], &URL{Address: e.Request.AbsoluteURL(link), Scraper: &scraper})
	})

	// Get all manga titles
	mangaCollector.OnHTML(source.MangaTitle, func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.Text)
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		if e.Index < len(scraper.manga[path]) {
			scraper.manga[path][e.Index].Info = title
		}
	})

	_ = mangaCollector.Limit(&colly.LimitRule{
		Parallelism: common.Parallelism,
		RandomDelay: time.Duration(source.RandomDelayMs) * time.Millisecond,
		DomainGlob:  "*",
	})

	// Paths collector
	chaptersCollector := collector.Clone()
	chaptersCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", common.Referer)
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
			u := &URL{
				Address:  e.Request.AbsoluteURL(link),
				Scraper:  &scraper,
				Index:    e.Index,
				Relation: scraper.chaptersCollector.Ctx,
			}

			urls = append(urls, u)
			scraper.chapters[path] = append(scraper.chapters[path], u)
		})

		urlsLength := len(urls)

		// Get all chapter titles
		html.ForEachWithBreak(source.ChapterTitle, func(i int, e *colly.HTMLElement) bool {
			title := strings.TrimSpace(e.Text)
			path := e.Request.AbsoluteURL(e.Request.URL.Path)

			if e.Index >= len(scraper.chapters[path]) {
				return false
			}

			scraper.chapters[path][e.Index].Info = title
			if source.ChaptersReversed {
				scraper.chapters[path][e.Index].Index = urlsLength - e.Index
			}
			return true
		})
	})
	_ = chaptersCollector.Limit(&colly.LimitRule{
		Parallelism: common.Parallelism,
		RandomDelay: time.Duration(source.RandomDelayMs) * time.Millisecond,
		DomainGlob:  "*",
	})

	// Pages collector
	pagesCollector := collector.Clone()
	pagesCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", common.Referer)
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

		// TODO: handle situations when ok is false
		attr, _ := util.Find(attributes, func(attr string) bool {
			return e.Attr(attr) != ""
		})

		link := e.Attr(attr)

		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		scraper.pages[path] = append(scraper.pages[path], &URL{Address: link, Scraper: &scraper, Index: e.Index, Relation: scraper.pagesCollector.Ctx})
	})
	_ = pagesCollector.Limit(&colly.LimitRule{
		Parallelism: common.Parallelism,
		RandomDelay: time.Duration(source.RandomDelayMs) * time.Millisecond,
		DomainGlob:  "*",
	})

	filesCollector := collector.Clone()
	filesCollector.CacheDir = ""
	filesCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", source.Base)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	})
	filesCollector.OnResponse(func(r *colly.Response) {
		scraper.Files.Set(r.Request.AbsoluteURL(r.Request.URL.Path), bytes.NewBuffer(r.Body))
	})

	scraper.mangaCollector = mangaCollector
	scraper.chaptersCollector = &ContextCollector[*URL]{
		Collector: chaptersCollector,
		Ctx:       nil,
	}
	scraper.pagesCollector = &ContextCollector[*URL]{
		Collector: pagesCollector,
		Ctx:       nil,
	}
	scraper.filesCollector = filesCollector

	return &scraper
}

// SearchManga searches for manga by name
func (s *Scraper) SearchManga(title string) ([]*URL, error) {
	// lowercase titles will produce the same results but will be useful for caching
	query := strings.ReplaceAll(title, " ", s.Source.WhitespaceEscape)
	address := fmt.Sprintf(s.Source.SearchTemplate, url.QueryEscape(strings.TrimSpace(strings.ToLower(query))))

	if urls, ok := s.manga[address]; ok {
		return urls, nil
	}

	err := s.mangaCollector.Visit(address)

	if err != nil {
		return nil, err
	}

	s.mangaCollector.Wait()
	return s.manga[address], nil
}

// GetChapters returns manga chapters for given manga url
func (s *Scraper) GetChapters(manga *URL) ([]*URL, error) {
	if urls, ok := s.chapters[manga.Address]; ok {
		return urls, nil
	}

	s.chaptersCollector.Ctx = manga
	err := s.chaptersCollector.Collector.Visit(manga.Address)
	if err != nil {
		return nil, err
	}

	s.chaptersCollector.Collector.Wait()
	s.chaptersCollector.Ctx = nil

	return s.chapters[manga.Address], nil
}

// GetPages returns manga pages for given chapter url
func (s *Scraper) GetPages(chapter *URL) ([]*URL, error) {
	if urls, ok := s.pages[chapter.Address]; ok {
		return urls, nil
	}

	s.pagesCollector.Ctx = chapter
	err := s.pagesCollector.Collector.Visit(chapter.Address)

	if err != nil {
		return nil, err
	}

	s.pagesCollector.Collector.Wait()
	s.pagesCollector.Ctx = nil

	return s.pages[chapter.Address], nil
}

// GetFile returns manga file for given page url
func (s *Scraper) GetFile(file *URL) (*bytes.Buffer, error) {
	if data, ok := s.Files.Get(file.Address); ok {
		return data, nil
	}

	err := s.filesCollector.Visit(file.Address)

	if err != nil {
		return nil, err
	}

	s.filesCollector.Wait()

	data, ok := s.Files.Get(file.Address)

	if ok {
		return data, nil
	}

	return nil, errors.New("Couldn't get file at " + file.Address)
}

// ResetFiles resets scraper files cache
func (s *Scraper) ResetFiles() {
	s.Files.Reset()
}
