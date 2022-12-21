package anilist

import (
	"github.com/metafates/mangal/key"
	"github.com/spf13/viper"
)

type Anilist struct {
	token string
}

// New cereates a new Anilist integration instance
func New() *Anilist {
	return &Anilist{}
}

func (a *Anilist) id() string {
	return viper.GetString(key.AnilistID)
}

func (a *Anilist) secret() string {
	return viper.GetString(key.AnilistSecret)
}

func (a *Anilist) code() string {
	return viper.GetString(key.AnilistCode)
}

// AuthURL returns the URL to authenticate with Anilist
func (a *Anilist) AuthURL() string {
	return "https://anilist.co/api/v2/oauth/authorize?client_id=" + a.id() + "&response_type=code&redirect_uri=https://anilist.co/api/v2/oauth/pin"
}
