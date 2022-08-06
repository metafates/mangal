package html

import (
	"github.com/PuerkitoBio/goquery"
	lua "github.com/yuin/gopher-lua"
)

const SelectionTypename = "selection"

var selectionMethods = map[string]lua.LGFunction{
	"find":     selectionFind,
	"each":     selectionEach,
	"attr":     selectionAttr,
	"first":    selectionFirst,
	"parent":   selectionParent,
	"text":     selectionText,
	"html":     selectionHtml,
	"hasClass": selectionHasClass,
	"is":       selectionIs,
	"next":     selectionNext,
	"prev":     selectionPrev,
}

func registerSelectionType(L *lua.LState) {
	mt := L.NewTypeMetatable(SelectionTypename)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), selectionMethods))
}

func checkSelection(L *lua.LState) *goquery.Selection {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*goquery.Selection); ok {
		return v
	}
	L.ArgError(1, "selection expected")
	return nil
}

func pushSelection(L *lua.LState, s *goquery.Selection) {
	ud := L.NewUserData()
	ud.Value = s
	L.SetMetatable(ud, L.GetTypeMetatable(SelectionTypename))
	L.Push(ud)
}

func selectionFind(L *lua.LState) int {
	s := checkSelection(L)
	selector := L.ToString(2)
	newS := s.Find(selector)
	pushSelection(L, newS)
	return 1
}

func selectionEach(L *lua.LState) int {
	s := checkSelection(L)
	fn := L.ToFunction(2)
	newS := s.Each(func(i int, s *goquery.Selection) {
		L.Push(fn)
		L.Push(lua.LNumber(i))
		pushSelection(L, s)
		if err := L.PCall(2, lua.MultRet, nil); err != nil {
			L.Error(lua.LString(err.Error()), 0)
			return
		}
	})
	pushSelection(L, newS)
	return 1
}

func selectionText(L *lua.LState) int {
	s := checkSelection(L)
	L.Push(lua.LString(s.Text()))
	return 1
}

func selectionHtml(L *lua.LState) int {
	s := checkSelection(L)
	html, err := s.Html()
	if err != nil {
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}
	L.Push(lua.LString(html))
	return 1
}

func selectionFirst(L *lua.LState) int {
	s := checkSelection(L)
	pushSelection(L, s.First())
	return 1
}

func selectionParent(L *lua.LState) int {
	s := checkSelection(L)
	pushSelection(L, s.Parent())
	return 1
}

func selectionAttr(L *lua.LState) int {
	s := checkSelection(L)
	attrName := L.ToString(2)
	attr, exists := s.Attr(attrName)
	L.Push(lua.LString(attr))
	L.Push(lua.LBool(exists))
	return 2
}

func selectionHasClass(L *lua.LState) int {
	s := checkSelection(L)
	className := L.ToString(2)
	L.Push(lua.LBool(s.HasClass(className)))
	return 1
}

func selectionIs(L *lua.LState) int {
	s := checkSelection(L)
	selector := L.ToString(2)
	L.Push(lua.LBool(s.Is(selector)))
	return 1
}

func selectionNext(L *lua.LState) int {
	s := checkSelection(L)
	pushSelection(L, s.Next())
	return 1
}

func selectionPrev(L *lua.LState) int {
	s := checkSelection(L)
	pushSelection(L, s.Prev())
	return 1
}
