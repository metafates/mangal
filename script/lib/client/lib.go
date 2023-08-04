package client

import (
	luadoc "github.com/mangalorg/gopher-luadoc"
	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/script/lib/util"
	lua "github.com/yuin/gopher-lua"
)

const (
	libName = "client"

	mangaTypeName         = libName + "_manga"
	volumeTypeName        = libName + "_volume"
	chapterTypeName       = libName + "_chapter"
	pageTypeName          = libName + "_page"
	pageWithImageTypeName = libName + "_page_with_image"
)

func Lib(client *libmangal.Client) *luadoc.Lib {
	classManga := &luadoc.Class{
		Name:        mangaTypeName,
		Description: "",
		Methods: []*luadoc.Method{
			{
				Name:        "info",
				Description: "",
				Value:       mangaInfo,
				Returns: []*luadoc.Param{
					{
						Name:        "info",
						Description: "",
						Type: luadoc.TableLiteral(
							"title", luadoc.String,
							"url", luadoc.String,
							"id", luadoc.String,
							"banner", luadoc.String,
							"cover", luadoc.String,
							"anilist_search", luadoc.String,
						),
					},
				},
			},
		},
	}

	classVolume := &luadoc.Class{
		Name:        volumeTypeName,
		Description: "",
		Methods: []*luadoc.Method{
			{
				Name:        "info",
				Description: "",
				Value:       volumeInfo,
				Returns: []*luadoc.Param{
					{
						Name:        "info",
						Description: "",
						Type: luadoc.TableLiteral(
							"number", luadoc.Number,
						),
					},
				},
			},
		},
	}

	classChapter := &luadoc.Class{
		Name:        chapterTypeName,
		Description: "",
		Methods: []*luadoc.Method{
			{
				Name:        "info",
				Description: "",
				Value:       chapterInfo,
				Returns: []*luadoc.Param{
					{
						Name:        "info",
						Description: "",
						Type: luadoc.TableLiteral(
							"title", luadoc.String,
							"number", luadoc.Number,
							"url", luadoc.String,
						),
					},
				},
			},
		},
	}

	classPage := &luadoc.Class{
		Name:        pageTypeName,
		Description: "Single page of the chapter",
	}

	classPageWithImage := &luadoc.Class{
		Name:        pageWithImageTypeName,
		Description: "Page with downloaded image",
		Methods: []*luadoc.Method{
			{
				Name:        "image",
				Description: "Get image bytes",
				Value:       pageImage,
				Returns: []*luadoc.Param{
					{
						Name:        "image",
						Description: "Image bytes",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "extension",
				Description: "GetExtension gets the image extension of this page. An extension must start with the dot. For example: .jpeg .png",
				Value:       pageExtension,
				Returns: []*luadoc.Param{
					{
						Name:        "extension",
						Description: "Image extension",
						Type:        luadoc.String,
					},
				},
			},
		},
	}

	return &luadoc.Lib{
		Name:        libName,
		Description: "Core mangal client (provider wrapper) functionality",
		Classes: []*luadoc.Class{
			classManga,
			classVolume,
			classChapter,
			classPageWithImage,
		},
		Funcs: []*luadoc.Func{
			{
				Name:        "search_mangas",
				Description: "Searches for mangas with the given query",
				Value:       newSearchMangas(client),
				Params: []*luadoc.Param{
					{
						Name: "query",
						Type: luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name: "mangas",
						Type: luadoc.List(classManga.Name),
					},
				},
			},
			{

				Name:        "manga_volumes",
				Description: "",
				Value:       newMangaVolumes(client),
				Params: []*luadoc.Param{
					{
						Name: "manga",
						Type: classManga.Name,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name: "volumes",
						Type: luadoc.List(classVolume.Name),
					},
				},
			},
			{

				Name:        "volume_chapters",
				Description: "",
				Value:       newVolumeChapters(client),
				Params: []*luadoc.Param{
					{
						Name: "volume",
						Type: classVolume.Name,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name: "chapters",
						Type: luadoc.List(classChapter.Name),
					},
				},
			},
			{
				Name:        "chapter_pages",
				Description: "Get chapter pages",
				Value:       newChapterPages(client),
				Params: []*luadoc.Param{
					{
						Name: "chapter",
						Type: classChapter.Name,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name: "pages",
						Type: luadoc.List(classPage.Name),
					},
				},
			},
			{
				Name:        "download_page",
				Description: "Download page",
				Value:       newDownloadPage(client),
				Params: []*luadoc.Param{
					{
						Name: "page",
						Type: classPage.Name,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name: "page_with_image",
						Type: classPageWithImage.Name,
					},
				},
			},
		},
	}
}

func mangaInfo(state *lua.LState) int {
	manga := util.Check[libmangal.Manga](state, 1)
	info := manga.Info()

	table := state.NewTable()

	table.RawSetString("title", lua.LString(info.Title))
	table.RawSetString("url", lua.LString(info.URL))
	table.RawSetString("id", lua.LString(info.ID))
	table.RawSetString("banner", lua.LString(info.Banner))
	table.RawSetString("cover", lua.LString(info.Cover))
	table.RawSetString("anilist_search", lua.LString(info.AnilistSearch))

	state.Push(table)
	return 1
}

func volumeInfo(state *lua.LState) int {
	volume := util.Check[libmangal.Volume](state, 1)
	info := volume.Info()

	table := state.NewTable()

	table.RawSetString("number", lua.LNumber(info.Number))

	state.Push(table)
	return 1
}

func chapterInfo(state *lua.LState) int {
	chapter := util.Check[libmangal.Chapter](state, 1)
	info := chapter.Info()

	table := state.NewTable()

	table.RawSetString("title", lua.LString(info.Title))
	table.RawSetString("url", lua.LString(info.URL))
	table.RawSetString("number", lua.LNumber(info.Number))

	state.Push(table)
	return 1
}

func missingClientError(state *lua.LState) int {
	state.RaiseError("missing provider client")
	return 0
}

func newSearchMangas(client *libmangal.Client) lua.LGFunction {
	if client == nil {
		return missingClientError
	}

	return func(state *lua.LState) int {
		query := state.CheckString(1)

		mangas, err := client.SearchMangas(state.Context(), query)
		util.Must(state, err)

		table := util.SliceToTable(state, mangas, func(manga libmangal.Manga) lua.LValue {
			return util.NewUserData(state, manga, mangaTypeName)
		})

		state.Push(table)
		return 1
	}
}

func newMangaVolumes(client *libmangal.Client) lua.LGFunction {
	if client == nil {
		return missingClientError
	}

	return func(state *lua.LState) int {
		manga := util.Check[libmangal.Manga](state, 1)

		volumes, err := client.MangaVolumes(state.Context(), manga)
		util.Must(state, err)

		table := util.SliceToTable(state, volumes, func(volume libmangal.Volume) lua.LValue {
			return util.NewUserData(state, volume, volumeTypeName)
		})

		state.Push(table)
		return 1
	}
}

func newVolumeChapters(client *libmangal.Client) lua.LGFunction {
	if client == nil {
		return missingClientError
	}

	return func(state *lua.LState) int {
		volume := util.Check[libmangal.Volume](state, 1)

		chapters, err := client.VolumeChapters(state.Context(), volume)
		util.Must(state, err)

		table := util.SliceToTable(state, chapters, func(chapter libmangal.Chapter) lua.LValue {
			return util.NewUserData(state, chapter, chapterTypeName)
		})

		state.Push(table)
		return 1
	}
}

func newChapterPages(client *libmangal.Client) lua.LGFunction {
	if client == nil {
		return missingClientError
	}

	return func(state *lua.LState) int {
		chapter := util.Check[libmangal.Chapter](state, 1)

		pages, err := client.ChapterPages(state.Context(), chapter)
		util.Must(state, err)

		table := util.SliceToTable(state, pages, func(page libmangal.Page) lua.LValue {
			return util.NewUserData(state, page, pageTypeName)
		})

		state.Push(table)
		return 1
	}
}

func newDownloadPage(client *libmangal.Client) lua.LGFunction {
	if client == nil {
		return missingClientError
	}

	return func(state *lua.LState) int {
		page := util.Check[libmangal.Page](state, 1)

		pageWithImage, err := client.DownloadPage(state.Context(), page)
		util.Must(state, err)

		util.Push(state, pageWithImage, pageWithImageTypeName)
		return 1
	}
}

func pageImage(state *lua.LState) int {
	page := util.Check[libmangal.PageWithImage](state, 1)
	image := page.GetImage()

	state.Push(lua.LString(image))
	return 1
}

func pageExtension(state *lua.LState) int {
	page := util.Check[libmangal.PageWithImage](state, 1)
	extension := page.GetExtension()

	state.Push(lua.LString(extension))
	return 1
}

func newDownloadChapter(client *libmangal.Client) lua.LGFunction {
	if client == nil {
		return missingClientError
	}

	return func(state *lua.LState) int {
		chapter := util.Check[libmangal.Chapter](state, 1)

		path, err := client.DownloadChapter(state.Context(), chapter, libmangal.DefaultDownloadOptions())
		util.Must(state, err)

		state.Push(lua.LString(path))
		return 1
	}
}
