package http

// Package http_client implements golang package http_client utility functionality for lua.

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds http_client to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local http_client = require("http_client")
func Preload(L *lua.LState) {
	L.PreloadModule("http_client", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	http_client_ud := L.NewTypeMetatable(`http_client_ud`)
	L.SetGlobal(`http_client_ud`, http_client_ud)
	L.SetField(http_client_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"do_request": DoRequest,
	}))

	http_request_ud := L.NewTypeMetatable(`http_request_ud`)
	L.SetGlobal(luaRequestType, http_request_ud)
	L.SetField(http_request_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"set_basic_auth": SetBasicAuth,
		"header_set":     HeaderSet,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"client":       New,
	"request":      NewRequest,
	"file_request": NewFileRequest,
}
