package mangapill

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/metafates/mangal/provider/generic"
	"net/url"
	"strings"
	"time"
)

var Config = &generic.Configuration{
	Name:            "Mangapill",
	Delay:           50 * time.Millisecond,
	Parallelism:     50,
	ReverseChapters: true,
	BaseURL:         "https://mangapill.com",
	GenerateSearchURL: func(query string) string {
		query = strings.ReplaceAll(query, " ", "+")
		query = strings.ToLower(query)
		query = strings.TrimSpace(query)
		template := "https://mangapill.com/search?q=%s&type=&status="
		return fmt.Sprintf(template, url.QueryEscape(query))
	},
	MangaExtractor: &generic.Extractor{
		Selector: "body > div.container.py-3 > div.my-3.grid.justify-end.gap-3.grid-cols-2.md\\:grid-cols-3.lg\\:grid-cols-5 > div",
		Name: func(selection *goquery.Selection) string {
			return strings.TrimSpace(selection.Find("div a div.leading-tight").Text())
		},
		URL: func(selection *goquery.Selection) string {
			return selection.Find("div a:first-child").AttrOr("href", "")
		},
		Volume: func(selection *goquery.Selection) string {
			return ""
		},
		Cover: func(selection *goquery.Selection) string {
			return selection.Find("img").AttrOr("data-src", "")
		},
	},
	ChapterExtractor: &generic.Extractor{
		Selector: "div[data-filter-list] a",
		Name: func(selection *goquery.Selection) string {
			return strings.TrimSpace(selection.Text())
		},
		URL: func(selection *goquery.Selection) string {
			return selection.AttrOr("href", "")
		},
		Volume: func(selection *goquery.Selection) string {
			return ""
		},
	},
	PageExtractor: &generic.Extractor{
		Selector: "picture img",
		URL: func(selection *goquery.Selection) string {
			return selection.AttrOr("data-src", "")
		},
	},
}
