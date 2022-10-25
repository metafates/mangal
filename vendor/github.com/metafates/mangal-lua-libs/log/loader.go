package log

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds log to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local log = require("log")
func Preload(L *lua.LState) {
	L.PreloadModule("log", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	loggerUD := L.NewTypeMetatable(`logger_ud`)
	L.SetGlobal(`logger_ud`, loggerUD)
	L.SetField(loggerUD, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"set_output": SetOutput,
		"set_prefix": SetPrefix,
		"set_flags":  SetFlags,
		"print":      Print,
		"printf":     Printf,
		"println":    Println,
		"close":      Close,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"new": New,
}
