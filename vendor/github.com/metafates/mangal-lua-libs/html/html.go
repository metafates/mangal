package html

import (
	lua "github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule("html", Loader())
}

func Loader() lua.LGFunction {
	var exports = map[string]lua.LGFunction{
		"parse": parseHTML(),
	}
	return func(L *lua.LState) int {
		mod := L.SetFuncs(L.NewTable(), exports)
		L.Push(mod)

		registerDocumentType(L)
		registerSelectionType(L)

		return 1
	}
}
