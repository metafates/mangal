package strings

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds strings to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local strings = require("strings")
func Preload(L *lua.LState) {
	L.PreloadModule("strings", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	registerStringsReader(L)
	registerStringsBuilder(L)

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"split":       Split,
	"trim":        Trim,
	"trim_space":  TrimSpace,
	"trim_prefix": TrimPrefix,
	"trim_suffix": TrimSuffix,
	"has_prefix":  HasPrefix,
	"has_suffix":  HasSuffix,
	"contains":    Contains,
	"new_reader":  newStringsReader,
	"new_builder": newStringsBuilder,
	"fields":      Fields,
}
