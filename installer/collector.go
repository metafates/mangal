package installer

import (
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"path/filepath"
)

var (
	user   = "metafates"
	repo   = "mangal-scrapers"
	branch = "main"
)

var collector = &githubFilesCollector{
	user:   user,
	repo:   repo,
	branch: branch,
}

// Scrapers gets available scrapers from GitHub repo.
// See https://github.com/metafates/mangal-scrapers
func Scrapers() ([]*Scraper, error) {
	err := collector.collect()
	if err != nil {
		return nil, err
	}

	return lo.FilterMap(collector.Files, func(f *GithubFile, _ int) (*Scraper, bool) {
		if filepath.Ext(f.Path) != ".lua" {
			return nil, false
		}

		return &Scraper{
			Name: util.FileStem(filepath.Base(f.Path)),
			URL:  f.Url,
		}, true
	}), nil
}
