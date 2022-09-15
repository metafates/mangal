package manganelo

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/metafates/mangal/provider/generic"
	"net/url"
	"strings"
	"time"
)

var Config = &generic.Configuration{
	Name:            "Manganelo",
	Delay:           50 * time.Millisecond,
	Parallelism:     50,
	ReverseChapters: true,
	BaseURL:         "https://ww5.manganelo.tv/",
	GenerateSearchURL: func(query string) string {
		query = strings.TrimSpace(query)
		query = strings.ToLower(query)
		query = url.QueryEscape(query)
		template := "https://ww5.manganelo.tv/search/%s"
		return fmt.Sprintf(template, query)
	},
	MangaExtractor: &generic.Extractor{
		Selector: ".search-story-item",
		Name: func(selection *goquery.Selection) string {
			return selection.Find("a.item-title").Text()
		},
		URL: func(selection *goquery.Selection) string {
			return selection.Find("a.item-title").AttrOr("href", "")
		},
		Cover: func(selection *goquery.Selection) string {
			return selection.Find(".item-img img").AttrOr("src", "")
		},
	},
	ChapterExtractor: &generic.Extractor{
		Selector: "li.a-h",
		Name: func(selection *goquery.Selection) string {
			name := selection.Find(".chapter-name").Text()
			if strings.HasPrefix(name, "Vol.") {
				splitted := strings.Split(name, " ")
				name = strings.Join(splitted[1:], " ")
			}
			return name
		},
		URL: func(selection *goquery.Selection) string {
			return selection.Find(".chapter-name").AttrOr("href", "")
		},
		Volume: func(selection *goquery.Selection) string {
			name := selection.Find(".chapter-name").Text()
			if strings.HasPrefix(name, "Vol.") {
				splitted := strings.Split(name, " ")
				return splitted[0]
			}
			return ""
		},
	},
	PageExtractor: &generic.Extractor{
		Selector: ".container-chapter-reader img",
		Name:     nil,
		URL: func(selection *goquery.Selection) string {
			return selection.AttrOr("data-src", "")
		},
	},
}
