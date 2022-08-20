package installer

import (
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"path/filepath"
)

var collector *githubFilesCollector

// Scrapers gets available scrapers from GitHub repo.
// See https://github.com/metafates/mangal-scrapers
func Scrapers() ([]*Scraper, error) {
	if collector == nil {
		setupCollector()
	}

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

func setupCollector() {
	collector = &githubFilesCollector{
		user:   viper.GetString(config.InstallerUser),
		repo:   viper.GetString(config.InstallerRepo),
		branch: viper.GetString(config.InstallerBranch),
	}
}
