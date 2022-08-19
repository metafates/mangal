package html

import (
	lua "github.com/yuin/gopher-lua"
)

type HTML struct{}

func New() *HTML {
	return &HTML{}
}

func (HTML) Name() string {
	return "html"
}

func (HTML) Loader() lua.LGFunction {
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
