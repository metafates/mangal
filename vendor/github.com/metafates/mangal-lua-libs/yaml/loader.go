package yaml

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds yaml to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local yaml = require("yaml")
func Preload(L *lua.LState) {
	L.PreloadModule("yaml", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	registerYAMLEncoder(L)
	registerYAMLDecoder(L)

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"decode":      Decode,
	"encode":      Encode,
	"new_encoder": newYAMLEncoder,
	"new_decoder": newYAMLDecoder,
}
