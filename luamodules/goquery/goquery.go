package goquery

import (
	lua "github.com/yuin/gopher-lua"
	"net/http"
)

type Goquery struct{}
type settings struct {
	client *http.Client
}

func New() *Goquery {
	return &Goquery{}
}

func (_ Goquery) Name() string {
	return "goquery"
}

func (_ Goquery) Loader() lua.LGFunction {
	var exports = map[string]lua.LGFunction{
		"doc": newDoc(),
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
