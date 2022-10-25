package json

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds json to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local json = require("json")
func Preload(L *lua.LState) {
	L.PreloadModule("json", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	registerJSONEncoder(L)
	registerDecoder(L)
	registerJsonDecodedObject(L)

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"tableIsObject": TableIsObject,
	"decode":        Decode,
	"encode":        Encode,
	"new_encoder":   newJSONEncoder,
	"new_decoder":   newJSONDecoder,
}
