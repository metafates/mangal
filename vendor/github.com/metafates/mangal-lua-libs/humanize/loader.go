package humanize

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds humanize to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local humanize = require("humanize")
func Preload(L *lua.LState) {
	L.PreloadModule("humanize", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"time":        Time,
	"parse_bytes": ParseBytes,
	"ibytes":      IBytes,
	"si":          SI,
}
