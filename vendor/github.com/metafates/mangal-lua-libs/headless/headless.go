package headless

import lua "github.com/yuin/gopher-lua"

func Preload(L *lua.LState) {
	L.PreloadModule("headless", Loader())
}

func Loader() lua.LGFunction {
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
