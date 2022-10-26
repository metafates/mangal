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
	limit   uint8 = 4
)

func normalizedName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func SetRelation(name string, to *Manga) error {
	err := relationCacher.Set(name, to.ID)
	if err != nil {
		return err
	}

	if _, ok := idCacher.Get(to.ID); !ok {
		return idCacher.Set(to.ID, to)
	}

	return nil
}

func GetRelation(name string) (*Manga, bool) {
	id, ok := relationCacher.Get(name)
	if !ok {
		return nil, false
	}

	manga, ok := idCacher.Get(id)
	return manga, ok
}

func FindClosest(name string) (*Manga, error) {
	if retries >= limit {
		retries = 0
		err := fmt.Errorf("no results found on Anilist for manga %s", name)
		log.Error(err)
		return nil, err
	}

	name = normalizedName(name)

	if id, ok := relationCacher.Get(name); ok {
		if manga, ok := idCacher.Get(id); ok {
			return manga, nil
		}
	}

	// search for manga on anilist
	mangas, err := SearchByName(name)
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
			normalizedName(a.Name()),
		) < levenshtein.Distance(
			name,
			normalizedName(b.Name()),
		)
	})

	log.Info("Found closest match: " + closest.Name())
	retries = 0

	if _, ok := relationCacher.Get(name); !ok {
		_ = relationCacher.Set(name, closest.ID)
	}
	_ = idCacher.Set(closest.ID, closest)
	return closest, nil
}
