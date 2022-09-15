package mangakakalot

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/metafates/mangal/provider/generic"
	"net/url"
	"strings"
	"time"
)

var Config = &generic.Configuration{
	Name:            "Mangakakalot",
	Delay:           50 * time.Millisecond,
	Parallelism:     50,
	ReverseChapters: true,
	BaseURL:         "https://mangakakalot.com/",
	GenerateSearchURL: func(query string) string {
		query = strings.ReplaceAll(query, " ", "_")
		query = strings.TrimSpace(query)
		query = strings.ToLower(query)
		query = url.QueryEscape(query)
		template := "https://mangakakalot.com/search/story/%s"
		return fmt.Sprintf(template, query)
	},
	MangaExtractor: &generic.Extractor{
		Selector: ".story-item",
		Name: func(selection *goquery.Selection) string {
			return selection.Find(".story-name a").Text()
		},
		URL: func(selection *goquery.Selection) string {
			return selection.Find(".story-name a").AttrOr("href", "")
		},
		Cover: func(selection *goquery.Selection) string {
			return selection.Find("img").AttrOr("src", "")
		},
	},
	ChapterExtractor: &generic.Extractor{
		Selector: "li.a-h",
		Name: func(selection *goquery.Selection) string {
			name := selection.Find("a").Text()
			if strings.HasPrefix(name, "Vol.") {
				splitted := strings.Split(name, " ")
				name = strings.Join(splitted[1:], " ")
			}
			return name
		},
		URL: func(selection *goquery.Selection) string {
			return selection.Find("a").AttrOr("href", "")
		},
		Volume: func(selection *goquery.Selection) string {
			name := selection.Find("a").Text()
			if strings.HasPrefix(name, "Vol.") {
				splitted := strings.Split(name, " ")
				return splitted[0]
			}
			return ""
		},
	},
	PageExtractor: &generic.Extractor{
		Selector: ".container-chapter-reader img",
		URL: func(selection *goquery.Selection) string {
			return selection.AttrOr("src", "")
		},
	},
}
