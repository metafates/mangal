package anilist

import (
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/fs"
	"github.com/mangalorg/mangal/path"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/bbolt"
	"github.com/philippgille/gokv/encoding"
	"log"
	"path/filepath"
)

var Client = newAnilist()

func newAnilist() *libmangal.Anilist {
	newPersistentStore := func(name string) (gokv.Store, error) {
		dir := filepath.Join(path.CacheDir(), "anilist")
		if err := fs.FS.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}

		return bbolt.NewStore(bbolt.Options{
			BucketName: name,
			Path:       filepath.Join(dir, name+".db"),
			Codec:      encoding.Gob,
		})
	}

	anilistOptions := libmangal.DefaultAnilistOptions()

	var err error
	anilistOptions.QueryToIDsStore, err = newPersistentStore("query-to-id")
	if err != nil {
		log.Fatal(err)
	}

	anilistOptions.IDToMangaStore, err = newPersistentStore("id-to-manga")
	if err != nil {
		log.Fatal(err)
	}

	anilistOptions.TitleToIDStore, err = newPersistentStore("title-to-id")
	if err != nil {
		log.Fatal(err)
	}

	anilistOptions.AccessTokenStore, err = newPersistentStore("access-token")
	if err != nil {
		log.Fatal(err)
	}

	anilist := libmangal.NewAnilist(anilistOptions)
	return &anilist
}
