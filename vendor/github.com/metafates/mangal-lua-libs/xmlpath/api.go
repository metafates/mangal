// Package xmlpath provides xmlpath for lua.
package xmlpath

import (
	"bytes"

	lua "github.com/yuin/gopher-lua"
	xmlpath "gopkg.in/xmlpath.v2"
)

type luaNode struct {
	*xmlpath.Node
}

type luaPath struct {
	*xmlpath.Path
}

type luaIter struct {
	*xmlpath.Iter
}

func checkLuaNode(L *lua.LState, n int) *luaNode {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaNode); ok {
		return v
	}
	L.ArgError(n, "xmlpath_node_ud expected")
	return nil
}

func newLuaNode(L *lua.LState, n *xmlpath.Node) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &luaNode{Node: n}
	L.SetMetatable(ud, L.GetTypeMetatable(`xmlpath_node_ud`))
	return ud
}

func checkLuaPath(L *lua.LState, n int) *luaPath {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaPath); ok {
		return v
	}
	L.ArgError(n, "xmlpath_path_ud expected")
	return nil
}

func newLuaPath(L *lua.LState, path *xmlpath.Path) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &luaPath{Path: path}
	L.SetMetatable(ud, L.GetTypeMetatable(`xmlpath_path_ud`))
	return ud
}

func checkLuaIter(L *lua.LState, n int) *luaIter {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaIter); ok {
		return v
	}
	L.ArgError(n, "xmlpath_iter_ud expected")
	return nil
}

// PathIter lua xmlpath_path_ud:iter() returns table of nodes
func PathIter(L *lua.LState) int {
	path := checkLuaPath(L, 1)
	node := checkLuaNode(L, 2)
	it := path.Path.Iter(node.Node)
	result := L.NewTable()
	i := 1
	for it.Next() {
		L.RawSetInt(result, i, newLuaNode(L, it.Node()))
		i++
	}
	L.Push(result)
	return 1
}

// IterNode lua xmlpath_iter_ud:node() returns node
func IterNode(L *lua.LState) int {
	iter := checkLuaIter(L, 1)
	L.Push(newLuaNode(L, iter.Iter.Node()))
	return 1
}

// NodeToString lua xmlpath_node_ud:string() returns string
func NodeToString(L *lua.LState) int {
	node := checkLuaNode(L, 1)
	L.Push(lua.LString(node.Node.String()))
	return 1
}

// Load lua xmlpath.load(xmlpath string) return (xmlpath_node_ud, err)
func Load(L *lua.LState) int {
	xmlpathStr := L.CheckString(1)
	r := bytes.NewReader([]byte(xmlpathStr))
	node, err := xmlpath.ParseHTML(r)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(newLuaNode(L, node))
	return 1
}

// Compile lua xmlpath.compile(xpath string) return (xmlpath_path_ud, err)
func Compile(L *lua.LState) int {
	xpathStr := L.CheckString(1)
	path, err := xmlpath.Compile(xpathStr)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(newLuaPath(L, path))
	return 1
}
