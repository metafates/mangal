package lib

import (
	luadoc "github.com/mangalorg/gopher-luadoc"
	"github.com/mangalorg/libmangal"
	luaprovidersdk "github.com/mangalorg/luaprovider/lib"
	"github.com/mangalorg/mangal/afs"
	"github.com/mangalorg/mangal/meta"
	"github.com/mangalorg/mangal/script/lib/client"
	"github.com/mangalorg/mangal/script/lib/json"
	"github.com/mangalorg/mangal/script/lib/prompt"
	lua "github.com/yuin/gopher-lua"
)

const libName = meta.AppName

type Options struct {
	Client  *libmangal.Client
	Anilist *libmangal.Anilist
}

func Lib(state *lua.LState, options Options) *luadoc.Lib {
	SDKOptions := luaprovidersdk.DefaultOptions()
	SDKOptions.FS = afs.Afero.Fs

	lib := &luadoc.Lib{
		Name:        libName,
		Description: meta.AppName + " scripting mode utilities",
		Libs: []*luadoc.Lib{
			luaprovidersdk.Lib(state, SDKOptions),
			prompt.Lib(),
			json.Lib(),
			client.Lib(options.Client),
		},
	}

	return lib
}

func Preload(state *lua.LState, options Options) {
	lib := Lib(state, options)
	state.PreloadModule(lib.Name, lib.Loader())
}
