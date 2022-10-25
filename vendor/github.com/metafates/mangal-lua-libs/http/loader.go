// Package http implements golang package http functionality for lua.

package http

import (
	client "github.com/metafates/mangal-lua-libs/http/client"
	util "github.com/metafates/mangal-lua-libs/http/util"
	lua "github.com/yuin/gopher-lua"
)

// Preload adds http to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//	local http = require("http")
func Preload(L *lua.LState) {
	L.PreloadModule("http", Loader)
	client.Preload(L)
	util.Preload(L)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	httpClientUd := L.NewTypeMetatable(`http_client_ud`)
	L.SetGlobal(`http_client_ud`, httpClientUd)
	L.SetField(httpClientUd, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"do_request": client.DoRequest,
	}))

	httpRequestUd := L.NewTypeMetatable(`http_request_ud`)
	L.SetGlobal(`http_request_ud`, httpRequestUd)
	L.SetField(httpRequestUd, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"set_basic_auth": client.SetBasicAuth,
		"header_set":     client.HeaderSet,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"client":         client.New,
	"request":        client.NewRequest,
	"file_request":   client.NewFileRequest,
	"query_escape":   util.QueryEscape,
	"query_unescape": util.QueryUnescape,
	"parse_url":      util.ParseURL,
	"build_url":      util.BuildURL,
}
