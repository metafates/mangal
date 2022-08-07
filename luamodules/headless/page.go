package headless

import (
	"github.com/go-rod/rod"
	lua "github.com/yuin/gopher-lua"
)

var pageMethods = map[string]lua.LGFunction{
	"waitLoad":             pageWaitLoad,
	"element":              pageElement,
	"elementR":             pageElementR,
	"elements":             pageElements,
	"elementByJS":          pageElementByJS,
	"waitElementsMoreThan": pageWaitElementsMoreThan,
	"navigate":             pageNavigate,
	"has":                  pageHas,
	"eval":                 pageEval,
	"html":                 pageHTML,
}

func registerPageType(L *lua.LState) {
	mt := L.NewTypeMetatable("browserPage")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), pageMethods))
}

func checkPage(L *lua.LState) *rod.Page {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*rod.Page); ok {
		return v
	}
	L.ArgError(1, "browserPage expected")
	return nil
}

func pageWaitLoad(L *lua.LState) int {
	page := checkPage(L)
	page.MustWaitLoad()

	return 0
}

func pageElement(L *lua.LState) int {
	p := checkPage(L)
	selector := L.CheckString(2)

	ud := L.NewUserData()
	ud.Value = p.MustElement(selector)
	L.SetMetatable(ud, L.GetTypeMetatable("pageElement"))

	L.Push(ud)
	return 1
}

func pageNavigate(L *lua.LState) int {
	p := checkPage(L)
	url := L.CheckString(2)
	p.MustNavigate(url)

	return 0
}

func pageElementByJS(L *lua.LState) int {
	p := checkPage(L)
	js := L.CheckString(2)

	ud := L.NewUserData()
	ud.Value = p.MustElementByJS(js)
	L.SetMetatable(ud, L.GetTypeMetatable("pageElement"))

	L.Push(ud)
	return 1
}

func pageElements(L *lua.LState) int {
	p := checkPage(L)
	selector := L.CheckString(2)
	els := p.MustElements(selector)

	table := L.NewTable()
	for i, el := range els {
		ud := L.NewUserData()
		ud.Value = el
		L.SetMetatable(ud, L.GetTypeMetatable("pageElement"))
		table.RawSetInt(i+1, ud)
	}

	L.Push(table)
	return 1
}

func pageHas(L *lua.LState) int {
	p := checkPage(L)
	selector := L.CheckString(2)
	L.Push(lua.LBool(p.MustHas(selector)))
	return 1
}

func pageEval(L *lua.LState) int {
	p := checkPage(L)
	js := L.CheckString(2)
	result := p.MustEval(js).Str()
	L.Push(lua.LString(result))
	return 1
}

func pageHTML(L *lua.LState) int {
	p := checkPage(L)
	html := p.MustHTML()
	L.Push(lua.LString(html))
	return 1
}

func pageElementR(L *lua.LState) int {
	p := checkPage(L)
	selector := L.CheckString(2)
	re := L.CheckString(3)

	ud := L.NewUserData()
	ud.Value = p.MustElementR(selector, re)
	L.SetMetatable(ud, L.GetTypeMetatable("pageElement"))

	L.Push(ud)
	return 1
}

func pageWaitElementsMoreThan(L *lua.LState) int {
	p := checkPage(L)
	selector := L.CheckString(2)
	count := L.CheckInt(3)
	p.MustWaitElementsMoreThan(selector, count)

	return 0
}
