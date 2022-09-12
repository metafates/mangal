package anilist

import (
	"encoding/json"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/log"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
)

var cache = anilistCache{
	Mem: make(map[string][]*Manga),
}

func init() {
	_ = cache.Init()
}

type anilistCache struct {
	Mem  map[string][]*Manga `json:"mangas"`
	file afero.File
}

func (a *anilistCache) Init() error {
	if a.file != nil {
		return nil
	}

	log.Info("Initializing anilist cacher")

	path := filepath.Join(os.TempDir(), constant.TempPrefix+"anilist_cache.json")
	exists, err := filesystem.Get().Exists(path)
	if err != nil {
		log.Warn(err)
		return err
	}

	if !exists {
		a.file, err = filesystem.Get().Create(path)
	} else {
		a.file, err = filesystem.Get().Open(path)
	}

	if err != nil {
		log.Warn(err)
		return err
	}

	var contents []byte
	_, err = a.file.Read(contents)
	if err != nil {
		log.Warn(err)
		return err
	}

	if len(contents) == 0 {
		return nil
	}

	var temp *anilistCache
	err = json.Unmarshal(contents, temp)
	if err != nil {
		log.Warn(err)
		return err
	}

	a.Mem = temp.Mem
	return nil
}

func (a *anilistCache) Get(name string) ([]*Manga, bool) {
	mangas, ok := a.Mem[name]
	if ok {
		log.Info("Found cached data in anilist cacher for " + name)
	}
	return mangas, ok
}

func (a *anilistCache) Set(name string, mangas []*Manga) error {
	log.Info("Setting anilist cacher entry")
	a.Mem[name] = mangas
	marshalled, err := json.Marshal(a)
	if err != nil {
		log.Warn(errors.Wrap(err, "anilist cache"))
		return err
	}

	_, err = a.file.Write(marshalled)
	return err
}
