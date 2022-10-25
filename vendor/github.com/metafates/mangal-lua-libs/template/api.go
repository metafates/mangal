// Package template implements template functionality for lua.
package template

import (
	"fmt"
	"io/ioutil"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

var (
	allEngines     = make(map[string]luaTemplateEngine, 0)
	allEnginesLock = &sync.Mutex{}
)

type luaTemplateEngine interface {
	Render(string, *lua.LTable) (string, error)
}

// RegisterTemplateEngine register template engine
func RegisterTemplateEngine(driver string, i luaTemplateEngine) {
	allEnginesLock.Lock()
	defer allEnginesLock.Unlock()

	allEngines[driver] = i
}

func checkEngine(L *lua.LState, n int) luaTemplateEngine {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(luaTemplateEngine); ok {
		return v
	}
	L.ArgError(n, "template expected")
	return nil
}

// Choose lua template.get(engine) returns (template_ud, err)
func Choose(L *lua.LState) int {
	allEnginesLock.Lock()
	defer allEnginesLock.Unlock()

	engine := L.CheckString(1)
	e, ok := allEngines[engine]
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("unknown template engine: %s", engine)))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = e
	L.SetMetatable(ud, L.GetTypeMetatable(`template_ud`))
	L.Push(ud)
	return 1
}

// Render lua template_ud:render(string, values) returns (string, err)
func Render(L *lua.LState) int {
	t := checkEngine(L, 1)
	body := L.CheckString(2)
	context := L.CheckTable(3)
	result, err := t.Render(body, context)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	return 1
}

// RenderFile lua template_ud:render(string, values) returns (string, err)
func RenderFile(L *lua.LState) int {
	t := checkEngine(L, 1)
	file := L.CheckString(2)
	context := L.CheckTable(3)
	body, err := ioutil.ReadFile(file)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result, err := t.Render(string(body), context)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	return 1
}
