package time

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds time to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local time = require("time")
func Preload(L *lua.LState) {
	L.PreloadModule("time", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"unix":      Unix,
	"unix_nano": UnixNano,
	"sleep":     Sleep,
	"parse":     Parse,
	"format":    Format,
}
