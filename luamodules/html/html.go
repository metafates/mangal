package html

import (
	lua "github.com/yuin/gopher-lua"
	"net/http"
)

type HTML struct{}
type settings struct {
	client *http.Client
}

func New() *HTML {
	return &HTML{}
}

func (_ HTML) Name() string {
	return "html"
}

func (_ HTML) Loader() lua.LGFunction {
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

type Option func(*settings)
