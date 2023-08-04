package anilist

import (
	luadoc "github.com/mangalorg/gopher-luadoc"
	"github.com/mangalorg/libmangal"
)

const libName = "anilist"

func Lib(anilist *libmangal.Anilist) *luadoc.Lib {
	return &luadoc.Lib{
		Name:        libName,
		Description: "Anilist operations",
	}
}
