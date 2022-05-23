package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Scraper struct {
	Source            *Source
	MangaCollector    *colly.Collector
	ChaptersCollector *colly.Collector
	PagesCollector    *colly.Collector
	FilesCollector    *colly.Collector

	Manga    sync.Map
	Chapters sync.Map
	Pages    sync.Map
	Files    sync.Map
}

var DefaultScraper = MakeSourceScraper(&DefaultSource)

// MakeSourceScraper returns new scraper with configures collectors
func MakeSourceScraper(source *Source) *Scraper {
	scraper := Scraper{
		Source: source,

		Manga:    sync.Map{},
		Chapters: sync.Map{},
		Pages:    sync.Map{},
		Files:    sync.Map{},
	}

	collector := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent(randomUserAgent()))

	_ = collector.Limit(&colly.LimitRule{
		Parallelism: Parallelism,
		RandomDelay: RandomDelay,
		DomainGlob:  fmt.Sprintf("*.%s/*", source.Base),
	})

	// Manga collector
	mangaCollector := collector.Clone()
	mangaCollector.OnHTML(source.MangaAnchor, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		m, _ := scraper.Manga.LoadOrStore(path, []*URL{})
		m = append(m.([]*URL), &URL{Address: e.Request.AbsoluteURL(link), Scraper: &scraper})
	})
	mangaCollector.OnHTML(source.MangaTitle, func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.Text)
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		m, _ := scraper.Manga.LoadOrStore(path, []*URL{})
		if e.Index < len(m.([]*URL)) {
			m.([]*URL)[e.Index].Info = title
		}
	})
	// Manga collector END

	// Chapters collector
	chaptersCollector := collector.Clone()
	chaptersCollector.OnHTML(source.ChapterAnchor, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		c, _ := scraper.Chapters.LoadOrStore(path, []*URL{})
		c = append(c.([]*URL), &URL{Address: e.Request.AbsoluteURL(link), Scraper: &scraper})
	})
	chaptersCollector.OnHTML(source.ChapterTitle, func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.Text)
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		c, _ := scraper.Chapters.LoadOrStore(path, []*URL{})
		if e.Index < len(c.([]*URL)) {
			c.([]*URL)[e.Index].Info = title
		}
	})
	// Chapters collector END

	// Pages collector
	pagesCollector := collector.Clone()
	pagesCollector.OnHTML(source.ReaderPage, func(e *colly.HTMLElement) {
		link := e.Attr("src")

		if link == "" {
			link = e.Attr("data-src")
		}

		filename := strconv.Itoa(e.Index)

		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		p, _ := scraper.Pages.LoadOrStore(path, []*URL{})
		p = append(p.([]*URL), &URL{Address: link, Info: filename, Scraper: &scraper})
	})
	// Pages collector END

	// File collector
	filesCollector := collector.Clone()
	filesCollector.OnResponse(func(r *colly.Response) {
		scraper.Files.Store(r.Request.AbsoluteURL(r.Request.URL.Path), &r.Body)
	})
	// File collector END

	scraper.MangaCollector = mangaCollector
	scraper.ChaptersCollector = chaptersCollector
	scraper.PagesCollector = pagesCollector
	scraper.FilesCollector = filesCollector

	return &scraper
}

func (s *Scraper) SearchManga(title string) ([]*URL, error) {
	address := fmt.Sprintf(s.Source.SearchTemplate, url.QueryEscape(title))

	if urls, ok := s.Manga.Load(address); ok {
		return urls.([]*URL), nil
	}

	err := s.MangaCollector.Visit(address)

	if err != nil {
		return nil, err
	}

	s.MangaCollector.Wait()
	m, _ := s.Manga.Load(address)
	return m.([]*URL), nil
}

func (s *Scraper) GetChapters(manga *URL) ([]*URL, error) {
	if urls, ok := s.Chapters.Load(manga.Address); ok {
		return urls.([]*URL), nil
	}

	err := s.ChaptersCollector.Visit(manga.Address)

	if err != nil {
		return nil, err
	}

	s.ChaptersCollector.Wait()

	// Add relation to this manga url for each chapter
	// It shouldn't affect performance since there won't be more than 1000 chapters as worst case
	chapters, _ := s.Chapters.Load(manga.Address)
	for _, chapter := range chapters.([]*URL) {
		chapter.Relation = manga
	}

	return chapters.([]*URL), nil
}

func (s *Scraper) GetPages(chapter *URL) ([]*URL, error) {
	if urls, ok := s.Pages.Load(chapter.Address); ok {
		return urls.([]*URL), nil
	}

	err := s.PagesCollector.Visit(chapter.Address)

	if err != nil {
		return nil, err
	}

	s.PagesCollector.Wait()
	p, _ := s.Pages.Load(chapter.Address)
	return p.([]*URL), nil
}

func (s *Scraper) GetFile(file *URL) (*[]byte, error) {
	if data, ok := s.Files.Load(file.Address); ok {
		return data.(*[]byte), nil
	}

	err := s.FilesCollector.Visit(file.Address)

	if err != nil {
		return nil, err
	}

	s.FilesCollector.Wait()

	data, _ := s.Files.Load(file.Address)

	return data.(*[]byte), nil
}

func (s *Scraper) CleanupFiles() {
	s.Files = sync.Map{}
}

func randomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	return strings.TrimSpace(userAgents[rand.Intn(len(userAgents))])
}

// USER AGENTS LIST BELOW

var userAgents = strings.Split(`
Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/37.0.2062.94 Chrome/37.0.2062.94 Safari/537.36
Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36
Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko
Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.0
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/600.8.9 (KHTML, like Gecko) Version/8.0.8 Safari/600.8.9
Mozilla/5.0 (iPad; CPU OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12H321 Safari/600.1.4
Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36
Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36
Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.10240
Mozilla/5.0 (Windows NT 6.3; WOW64; rv:40.0) Gecko/20100101 Firefox/40.0
Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko
Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36
Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko
Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/6.0)
Mozilla/5.0 (Windows NT 6.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36
Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36
Mozilla/5.0 (X11; CrOS x86_64 7077.134.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.156 Safari/537.36
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/600.7.12 (KHTML, like Gecko) Version/7.1.7 Safari/537.85.16
Mozilla/5.0 (Windows NT 6.0; rv:40.0) Gecko/20100101 Firefox/40.0
Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv:40.0) Gecko/20100101 Firefox/40.0
Mozilla/5.0 (iPad; CPU OS 8_1_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B466 Safari/600.1.4
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_2) AppleWebKit/600.3.18 (KHTML, like Gecko) Version/8.0.3 Safari/600.3.18
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36
Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36
Mozilla/5.0 (Windows NT 6.1; Win64; x64; Trident/7.0; rv:11.0) like Gecko
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36
Mozilla/5.0 (iPad; CPU OS 8_1_2 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B440 Safari/600.1.4
Mozilla/5.0 (Linux; U; Android 4.0.3; en-us; KFTT Build/IML74K) AppleWebKit/537.36 (KHTML, like Gecko) Silk/3.68 like Chrome/39.0.2171.93 Safari/537.36
Mozilla/5.0 (iPad; CPU OS 8_2 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12D508 Safari/600.1.4
Mozilla/5.0 (Windows NT 6.1; WOW64; rv:39.0) Gecko/20100101 Firefox/39.0
Mozilla/5.0 (iPad; CPU OS 7_1_1 like Mac OS X) AppleWebKit/537.51.2 (KHTML, like Gecko) Version/7.0 Mobile/11D201 Safari/9537.53
Mozilla/5.0 (Linux; U; Android 4.4.3; en-us; KFTHWI Build/KTU84M) AppleWebKit/537.36 (KHTML, like Gecko) Silk/3.68 like Chrome/39.0.2171.93 Safari/537.36
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/600.6.3 (KHTML, like Gecko) Version/7.1.6 Safari/537.85.15
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_2) AppleWebKit/600.4.10 (KHTML, like Gecko) Version/8.0.4 Safari/600.4.10
Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; rv:40.0) Gecko/20100101 Firefox/40.0
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.78.2 (KHTML, like Gecko) Version/7.0.6 Safari/537.78.2
`, "\n")
