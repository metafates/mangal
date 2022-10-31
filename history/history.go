package history

import (
	"github.com/metafates/mangal/cache"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/integration"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/where"
	"github.com/samber/mo"
	"github.com/spf13/viper"
	"time"
)

var cacher = cache.New[map[string]*SavedChapter](
	where.History(),
	&cache.Options{
		ExpireEvery: mo.None[time.Duration](),
	},
)

// Get returns all chapters from the history file
func Get() (chapters map[string]*SavedChapter, err error) {
	cached, ok := cacher.Get().Get()
	if !ok {
		return make(map[string]*SavedChapter), nil
	}

	return cached, nil
}

// Save saves the chapter to the history file
func Save(chapter *source.Chapter) error {
	if viper.GetBool(constant.AnilistEnable) {
		go func() {
			log.Info("Saving chapter to anilist")
			err := integration.Anilist.MarkRead(chapter)
			if err != nil {
				log.Warn("Saving chapter to anilist failed: " + err.Error())
			}
		}()
	}

	saved, err := Get()
	if err != nil {
		return err
	}

	savedChapter := newSavedChapter(chapter)
	saved[savedChapter.encode()] = savedChapter

	return cacher.Set(saved)
}

// Remove removes the chapter from the history file
func Remove(chapter *SavedChapter) error {
	saved, err := Get()
	if err != nil {
		return err
	}

	delete(saved, chapter.encode())

	return cacher.Set(saved)
}
