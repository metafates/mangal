package storage

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds storage to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local storage = require("storage")
func Preload(L *lua.LState) {
	L.PreloadModule("storage", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	storage_ud := L.NewTypeMetatable(`storage_ud`)
	L.SetGlobal(`storage_ud`, storage_ud)
	L.SetField(storage_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"get":   Get,
		"set":   Set,
		"sync":  Sync,
		"close": Close,
		"keys":  Keys,
		"dump":  Dump,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"open": New,
}
