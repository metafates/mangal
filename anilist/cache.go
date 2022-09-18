package anilist

import (
	"encoding/json"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/util"
	"github.com/metafates/mangal/where"
	"io"
	"os"
	"path/filepath"
)

var cache = anilistCache{
	data: &anilistCacheData{Mangas: make(map[string]*Manga)},
}

type anilistCacheData struct {
	Mangas map[string]*Manga `json:"mangas"`
}

type anilistCache struct {
	data        *anilistCacheData
	path        string
	initialized bool
}

func (a *anilistCache) Init() error {
	if a.initialized {
		return nil
	}

	log.Debug("Initializing anilist cacher")

	path := filepath.Join(where.Cache(), "anilist_cache.json")
	a.path = path
	log.Debugf("Opening anilist cache file at %s", path)
	file, err := filesystem.Api().OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)

	if err != nil {
		log.Warn(err)
		return err
	}

	defer util.Ignore(file.Close)

	contents, err := io.ReadAll(file)
	if err != nil {
		log.Warn(err)
		return err
	}

	if len(contents) == 0 {
		log.Debug("Anilist cache file is empty, skipping unmarshal")
		return nil
	}

	err = json.Unmarshal(contents, a.data)
	if err != nil {
		log.Warn(err)
		return err
	}

	log.Debugf("Anilist cache file unmarshalled successfully, len is %d", len(a.data.Mangas))
	return nil
}

func (a *anilistCache) Get(name string) (*Manga, bool) {
	_ = a.Init()

	mangas, ok := a.data.Mangas[normalizeName(name)]
	return mangas, ok
}

func (a *anilistCache) Set(name string, manga *Manga) error {
	_ = a.Init()

	log.Debug("Setting anilist cacher entry")
	a.data.Mangas[normalizeName(name)] = manga
	marshalled, err := json.Marshal(a.data)
	if err != nil {
		log.Warn(err)
		return err
	}

	file, err := filesystem.Api().OpenFile(a.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)

	_, err = file.Write(marshalled)
	if err != nil {
		log.Warn(err)
	}

	return err
}
