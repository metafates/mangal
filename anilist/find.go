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

	if id := idCacher.Get(to.ID); id.IsAbsent() {
		return idCacher.Set(to.ID, to)
	}

	return nil
}

func FindClosest(name string) (*Manga, error) {
	if retries >= limit {
		retries = 0
		err := fmt.Errorf("no results found on Anilist for manga %s", name)
		log.Error(err)
		return nil, err
	}

	name = normalizedName(name)
	id := relationCacher.Get(name)
	if id.IsPresent() {
		if manga := idCacher.Get(id.MustGet()); manga.IsPresent() {
			return manga.MustGet(), nil
		}
	}

	// search for manga on anilist
	mangas, err := SearchByName(name)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if id.IsPresent() {
		found, ok := lo.Find(mangas, func(item *Manga) bool {
			return item.ID == id.MustGet()
		})

		if ok {
			return found, nil
		}

		// there should be a manga with the id in the cache, but it wasn't found
		// this means that the manga was deleted from anilist
		// remove the id from the cache
		_ = relationCacher.Delete(name)
		log.Infof("Manga with id %d was deleted from Anilist", id.MustGet())
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

	if id := relationCacher.Get(name); id.IsAbsent() {
		_ = relationCacher.Set(name, closest.ID)
	}
	_ = idCacher.Set(closest.ID, closest)
	return closest, nil
}
