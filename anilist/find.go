package anilist

import (
	"fmt"
	levenshtein "github.com/ka-weihe/fast-levenshtein"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"strings"
)

var (
	retries uint8
	limit   uint8 = 3
)

func FindClosest(name string) (*Manga, error) {
	if retries >= limit {
		retries = 0
		err := fmt.Errorf("no results found on Anilist for manga %s", name)
		log.Error(err)
		return nil, err
	}

	//if cached, ok := cache[name]; ok {
	//	log.Info("Found cached manga: " + cached.Name())
	//	return cached, nil
	//}

	// search for manga on anilist
	urls, err := Search(name)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if len(urls) == 0 {
		// try again with a different name
		retries++
		words := strings.Split(name, " ")
		if len(words) == 1 {
			// trigger limit
			retries = limit
			return FindClosest("")
		}

		// one word less
		alternateName := strings.Join(words[:util.Max(len(words)-1, 1)], " ")
		log.Infof(`No results found on Anilist for manga "%s", trying "%s"`, name, alternateName)
		return FindClosest(alternateName)
	}

	// find the closest match
	closest := lo.MinBy(urls, func(a, b *Manga) bool {
		return levenshtein.Distance(name, a.Name()) < levenshtein.Distance(name, b.Name())
	})

	log.Info("Found closest match: " + closest.Name())
	//cache[name] = closest
	retries = 0
	return closest, nil
}
