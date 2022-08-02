package headless

import (
	"github.com/go-rod/rod"
	lua "github.com/yuin/gopher-lua"
)

var pageMethods = map[string]lua.LGFunction{
	"waitLoad":    waitPage,
	"element":     selectElement,
	"elements":    selectElements,
	"elementByJS": selectElementByJS,
	"navigate":    pageNavigate,
	"exists":      exists,
	"evalJS":      evalJS,
	"html":        pageHTML,
}

func registerPageType(L *lua.LState) {
	mt := L.NewTypeMetatable("page")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), pageMethods))
}

func checkPage(L *lua.LState) *rod.Page {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*rod.Page); ok {
		return v
	}
	L.ArgError(1, "page expected")
	return nil
}

func waitPage(L *lua.LState) int {
	page := checkPage(L)
	page.MustWaitLoad()

	return 0
}

func selectElement(L *lua.LState) int {
	page := checkPage(L)
	selector := L.ToString(2)
	element := page.MustElement(selector)

	ud := L.NewUserData()
	ud.Value = element
	L.SetMetatable(ud, L.GetTypeMetatable("element"))

	L.Push(ud)
	return 1
}

func pageNavigate(L *lua.LState) int {
	page := checkPage(L)
	url := L.ToString(2)
	page.MustNavigate(url)

	return 0
}

func selectElementByJS(L *lua.LState) int {
	page := checkPage(L)
	selector := L.ToString(2)
	element := page.MustElementByJS(selector)

	ud := L.NewUserData()
	ud.Value = element
	L.SetMetatable(ud, L.GetTypeMetatable("element"))

	L.Push(ud)
	return 1
}

func selectElements(L *lua.LState) int {
	page := checkPage(L)
	selector := L.ToString(2)
	elements := page.MustElements(selector)

	table := L.NewTable()
	for i, element := range elements {
		ud := L.NewUserData()
		ud.Value = element
		L.SetMetatable(ud, L.GetTypeMetatable("element"))
		table.RawSetInt(i+1, ud)
	}

	L.Push(table)
	return 1
}

func exists(L *lua.LState) int {
	page := checkPage(L)
	selector := L.ToString(2)
	has := page.MustHas(selector)
	L.Push(lua.LBool(has))
	return 1
}

func evalJS(L *lua.LState) int {
	page := checkPage(L)
	js := L.ToString(2)
	result := page.MustEval(js).Str()
	L.Push(lua.LString(result))
	return 1
}

func pageHTML(L *lua.LState) int {
	page := checkPage(L)
	html := page.MustHTML()
	L.Push(lua.LString(html))
	return 1
}
