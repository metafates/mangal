package headless

import (
	"github.com/go-rod/rod"
	lua "github.com/yuin/gopher-lua"
)

var browserMethods = map[string]lua.LGFunction{
	"open": browserOpenPage,
}

func newBrowser() lua.LGFunction {
	return func(L *lua.LState) int {
		browser := rod.New()
		err := browser.Connect()

		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		ud := L.NewUserData()
		ud.Value = browser
		L.SetMetatable(ud, L.GetTypeMetatable("browser"))

		L.Push(ud)
		L.Push(lua.LNil)
		return 2
	}
}

func registerBrowserType(L *lua.LState) {
	mt := L.NewTypeMetatable("browser")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), browserMethods))
}

func checkBrowser(L *lua.LState) *rod.Browser {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*rod.Browser); ok {
		return v
	}
	L.ArgError(1, "browser expected")
	return nil
}

func browserOpenPage(L *lua.LState) int {
	browser := checkBrowser(L)
	url := L.ToString(2)
	page := browser.MustPage(url)

	ud := L.NewUserData()
	ud.Value = page
	L.SetMetatable(ud, L.GetTypeMetatable("page"))

	L.Push(ud)
	return 1
}
