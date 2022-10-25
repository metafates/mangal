package goos

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds goos to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local goos = require("goos")
func Preload(L *lua.LState) {
	L.PreloadModule("goos", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"stat":         Stat,
	"hostname":     Hostname,
	"get_pagesize": Getpagesize,
	"mkdir_all":    MkdirAll,
}
