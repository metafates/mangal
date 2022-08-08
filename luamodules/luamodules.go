package luamodules

import (
	"github.com/metafates/mangal/luamodules/core"
	"github.com/metafates/mangal/luamodules/headless"
	"github.com/metafates/mangal/luamodules/html"
	lua "github.com/yuin/gopher-lua"
)

type LuaModule interface {
	Name() string
	Loader() lua.LGFunction
}

var modules = []LuaModule{
	html.New(),
	headless.New(),
}

func PreloadAll(L *lua.LState) {
	// special case for core module
	// because it loads modules slightly differently
	core.Preload(L)

	for _, module := range modules {
		L.PreloadModule(module.Name(), module.Loader())
	}
}
