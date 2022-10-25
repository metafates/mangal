// Package http_util implements golang package http utility functionality for lua.

package util

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds http to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local http_util = require("http_util")
func Preload(L *lua.LState) {
	L.PreloadModule("http_util", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"query_escape":   QueryEscape,
	"query_unescape": QueryUnescape,
	"parse_url":      ParseURL,
	"build_url":      BuildURL,
}
