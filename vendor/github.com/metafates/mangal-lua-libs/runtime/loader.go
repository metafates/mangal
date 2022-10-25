package runtime

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds runtime to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local runtime = require("runtime")
func Preload(L *lua.LState) {
	L.PreloadModule("runtime", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"goos":   GOOS,
	"goarch": GOARCH,
}
