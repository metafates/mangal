package anilist

import (
	"encoding/json"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/log"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var cache = anilistCache{
	Mem: make(map[string]*Manga),
}

type anilistCache struct {
	Mem  map[string]*Manga `json:"mangas"`
	file afero.File
}

func (a *anilistCache) Init() error {
	if a.file != nil {
		return nil
	}

	log.Debug("Initializing anilist cacher")

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Warn(err)
		return err
	}

	path := filepath.Join(cacheDir, constant.CachePrefix+"anilist_cache.json")
	log.Debugf("Opening anilist cache file at %s", path)
	a.file, err = filesystem.Get().OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		log.Warn(err)
		return err
	}

	contents, err := io.ReadAll(a.file)
	if err != nil {
		log.Warn(err)
		return err
	}

	if len(contents) == 0 {
		log.Debug("Anilist cache file is empty, skipping unmarshal")
		return nil
	}

	var temp anilistCache
	err = json.Unmarshal(contents, &temp)
	if err != nil {
		log.Warn(err)
		return err
	}

	log.Debugf("Anilist cache file unmarshalled successfully, len is %d", len(temp.Mem))
	a.Mem = temp.Mem
	return nil
}

func (a *anilistCache) Get(name string) (*Manga, bool) {
	if a.file == nil {
		_ = a.Init()
	}
	mangas, ok := a.Mem[a.formatName(name)]
	return mangas, ok
}

func (a *anilistCache) Set(name string, manga *Manga) error {
	log.Debug("Setting anilist cacher entry")
	a.Mem[a.formatName(name)] = manga
	marshalled, err := json.Marshal(a)
	if err != nil {
		log.Warn(err)
		return err
	}

	_, _ = a.file.Seek(0, 0)
	_, err = a.file.Write(marshalled)
	if err != nil {
		log.Warn(err)
	}
	return err
}

func (a *anilistCache) formatName(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}
