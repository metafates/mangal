package ioutil

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds ioutil to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local ioutil = require("ioutil")
func Preload(L *lua.LState) {
	L.PreloadModule("ioutil", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"read_file":  ReadFile,
	"write_file": WriteFile,
	"copy":       Copy,
	"copyn":      CopyN,
}
