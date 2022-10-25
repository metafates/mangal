package html

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
	lua "github.com/yuin/gopher-lua"
)

const DocumentTypename = "document"

var documentMethods = map[string]lua.LGFunction{
	"find": documentFind,
}

func registerDocumentType(L *lua.LState) {
	mt := L.NewTypeMetatable(DocumentTypename)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), documentMethods))
}

func parseHTML() lua.LGFunction {
	return func(L *lua.LState) int {
		docData := L.ToString(1)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(docData))
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		ud := L.NewUserData()
		ud.Value = doc
		L.SetMetatable(ud, L.GetTypeMetatable(DocumentTypename))
		L.Push(ud)
		L.Push(lua.LNil)
		return 2
	}
}

func checkDocument(L *lua.LState) *goquery.Document {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*goquery.Document); ok {
		return v
	}
	L.ArgError(1, "document expected")
	return nil
}

func documentFind(L *lua.LState) int {
	doc := checkDocument(L)
	selector := L.ToString(2)
	s := doc.Find(selector)
	pushSelection(L, s)
	return 1
}
