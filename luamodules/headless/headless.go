package headless

import lua "github.com/yuin/gopher-lua"

type Headless struct{}

func New() *Headless {
	return &Headless{}
}

func (_ Headless) Name() string {
	return "headless"
}

func (_ Headless) Loader() lua.LGFunction {
	var exports = map[string]lua.LGFunction{
		"browser": newBrowser(),
	}

	return func(L *lua.LState) int {
		mod := L.SetFuncs(L.NewTable(), exports)
		L.Push(mod)

		registerBrowserType(L)
		registerPageType(L)
		registerElementType(L)

		return 1
	}
}
