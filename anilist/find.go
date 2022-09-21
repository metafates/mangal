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

func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func SetRelation(name string, to *Manga) error {
	return cache.Set(name, to)
}

func GetRelation(name string) (*Manga, bool) {
	return cache.Get(name)
}

func FindClosest(name string) (*Manga, error) {
	if retries >= limit {
		retries = 0
		err := fmt.Errorf("no results found on Anilist for manga %s", name)
		log.Error(err)
		return nil, err
	}

	name = normalizeName(name)

	if manga, ok := cache.Get(name); ok {
		return manga, nil
	}

	// search for manga on anilist
	mangas, err := Search(name)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if len(mangas) == 0 {
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
	closest := lo.MinBy(mangas, func(a, b *Manga) bool {
		return levenshtein.Distance(
			name,
			normalizeName(a.Name()),
		) < levenshtein.Distance(
			name,
			normalizeName(b.Name()),
		)
	})

	log.Info("Found closest match: " + closest.Name())
	retries = 0
	_ = cache.Set(name, closest)
	return closest, nil
}
