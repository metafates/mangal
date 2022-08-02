package headless

import (
	"github.com/go-rod/rod"
	lua "github.com/yuin/gopher-lua"
)

var elementMethods = map[string]lua.LGFunction{
	"input":     elementInput,
	"click":     elementClick,
	"text":      elementText,
	"attribute": elementAttribute,
	"html":      elementHtml,
}

func registerElementType(L *lua.LState) {
	mt := L.NewTypeMetatable("element")
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), elementMethods))
}

func checkElement(L *lua.LState) *rod.Element {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*rod.Element); ok {
		return v
	}
	L.ArgError(1, "element expected")
	return nil
}

func elementInput(L *lua.LState) int {
	element := checkElement(L)
	value := L.ToString(2)
	element.MustInput(value)

	return 0
}

func elementClick(L *lua.LState) int {
	element := checkElement(L)
	element.MustClick()

	return 0
}

func elementText(L *lua.LState) int {
	element := checkElement(L)
	text := element.MustText()

	L.Push(lua.LString(text))
	return 1
}

func elementAttribute(L *lua.LState) int {
	element := checkElement(L)
	name := L.ToString(2)
	value := element.MustAttribute(name)

	L.Push(lua.LString(*value))
	return 1
}

func elementHtml(L *lua.LState) int {
	element := checkElement(L)
	html := element.MustHTML()

	L.Push(lua.LString(html))
	return 1
}
