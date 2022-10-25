// Package json implements json decode/encode functionality for lua.
// original code: https://github.com/layeh/gopher-json
package json

import (
	lua "github.com/yuin/gopher-lua"
)

// Decode lua json.decode(string) returns (table, err)
func Decode(L *lua.LState) int {
	str := L.CheckString(1)

	value, err := ValueDecode(L, []byte(str))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(value)
	return 1
}

// Encode lua json.encode(obj) returns (string, err)
func Encode(L *lua.LState) int {
	value := L.CheckAny(1)

	data, err := ValueEncode(value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(string(data)))
	return 1
}

//TableIsObject lua json.tableIsObject marks a table as an object (to distinguish between [] and {})
func TableIsObject(L *lua.LState) int {
	table := L.CheckTable(1)
	L.SetMetatable(table, L.GetTypeMetatable(jsonTableIsObject))
	return 0
}
