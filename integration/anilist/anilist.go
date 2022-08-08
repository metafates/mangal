package anilist

import (
	"github.com/metafates/mangal/config"
	"github.com/spf13/viper"
)

type Anilist struct {
	token string
	cache map[string]*anilistManga
}

func New() *Anilist {
	return &Anilist{
		cache: make(map[string]*anilistManga),
	}
}

func (a *Anilist) id() string {
	return viper.GetString(config.AnilistID)
}

func (a *Anilist) secret() string {
	return viper.GetString(config.AnilistSecret)
}

func (a *Anilist) code() string {
	return viper.GetString(config.AnilistCode)
}

func (a *Anilist) AuthURL() string {
	return "https://anilist.co/api/v2/oauth/authorize?client_id=" + a.id() + "&response_type=code&redirect_uri=https://anilist.co/api/v2/oauth/pin"
}
