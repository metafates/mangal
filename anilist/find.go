package anilist

import (
	levenshtein "github.com/ka-weihe/fast-levenshtein"
	"github.com/metafates/mangal/log"
	"github.com/samber/lo"
)

func FindClosest(name string) (*Manga, error) {
	if cached, ok := cache[name]; ok {
		log.Info("Found cached manga: " + cached.Name())
		return cached, nil
	}

	// search for manga on anilist
	urls, err := Search(name)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if len(urls) == 0 {
		log.Warn("No results found on Anilist for manga" + name)
		return nil, nil
	}

	// find the closest match
	closest := lo.MinBy(urls, func(a, b *Manga) bool {
		return levenshtein.Distance(name, a.Name()) < levenshtein.Distance(name, b.Name())
	})

	log.Info("Found closest match: " + closest.Name())
	cache[name] = closest
	return closest, nil
}
