package anilist

import (
	"context"

	luadoc "github.com/mangalorg/gopher-luadoc"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/script/lib/util"
	lua "github.com/yuin/gopher-lua"
)

const (
	libName = "anilist"

	mangaTypeName = libName + "_manga"
)

func Lib(anilist *libmangal.Anilist) *luadoc.Lib {
	return &luadoc.Lib{
		Name:        libName,
		Description: "Anilist operations",
		Funcs: []*luadoc.Func{
			{
				Name:        "search_mangas",
				Description: "Search mangas on Anilist",
				Value:       newSearchMangas(anilist),
				Params: []*luadoc.Param{
					{
						Name:        "query",
						Description: "Query to search",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "mangas",
						Description: "Anilist mangas",
						Type:        luadoc.List(mangaTypeName),
					},
				},
			},
			{
				Name:        "find_closest_mangas",
				Description: "",
				Value:       newFindClosestManga(anilist),
				Params: []*luadoc.Param{
					{
						Name:        "title",
						Description: "Manga title to search for",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "manga",
						Description: "Closest manga found",
						Type:        mangaTypeName,
					},
					{
						Name:        "found",
						Description: "Whether the closest manga was found",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "bind_title_with_id",
				Description: "Binds manga title to anilist manga id",
				Value:       newBindTitleWithID(anilist),
				Params: []*luadoc.Param{
					{
						Name:        "title",
						Description: "Manga title to use for binding",
						Type:        luadoc.String,
					},
					{
						Name:        "id",
						Description: "Anilist manga ID to bind to",
						Type:        luadoc.String,
					},
				},
			},
		},
	}
}

func missingAnilistError(state *lua.LState) int {
	state.RaiseError("anilist client is missing")
	return 0
}

func newSearchMangas(anilist *libmangal.Anilist) lua.LGFunction {
	if anilist == nil {
		return missingAnilistError
	}

	return func(state *lua.LState) int {
		query := state.CheckString(1)

		mangas, err := anilist.SearchMangas(state.Context(), query)
		util.Must(state, err)

		table := util.SliceToTable(state, mangas, func(manga libmangal.AnilistManga) lua.LValue {
			return util.NewUserData(state, manga, mangaTypeName)
		})

		state.Push(table)
		return 1
	}
}

func newFindClosestManga(anilist *libmangal.Anilist) lua.LGFunction {
	if anilist == nil {
		return missingAnilistError
	}

	return func(state *lua.LState) int {
		title := state.CheckString(1)

		manga, found, err := anilist.FindClosestManga(context.Background(), title)
		util.Must(state, err)

		util.Push(state, manga, mangaTypeName)
		state.Push(lua.LBool(found))
		return 2
	}
}

func newBindTitleWithID(anilist *libmangal.Anilist) lua.LGFunction {
	if anilist == nil {
		return missingAnilistError
	}

	return func(state *lua.LState) int {
		title := state.CheckString(1)
		ID := state.CheckInt(2)

		err := anilist.BindTitleWithID(title, ID)
		util.Must(state, err)

		return 0
	}
}
