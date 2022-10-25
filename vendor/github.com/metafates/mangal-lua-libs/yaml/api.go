// Package yaml implements yaml decode functionality for lua.
package yaml

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
	yaml "gopkg.in/yaml.v2"
)

// Decode lua yaml.decode(string) returns (table, error)
func Decode(L *lua.LState) int {
	str := L.CheckString(1)

	var value interface{}
	err := yaml.Unmarshal([]byte(str), &value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(fromYAML(L, value))
	return 1
}

// Encode lua yaml.encode(any) returns (string, error)
func Encode(L *lua.LState) int {
	arg := L.CheckAny(1)
	var value interface{}
	err := L.GPCall(func(L *lua.LState) int {
		visited := make(map[*lua.LTable]bool)
		value = toYAML(L, visited, arg)
		return 0
	}, lua.LNil)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	data, err := yaml.Marshal(value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(data))
	return 1
}

func tableIsSlice(table *lua.LTable) bool {
	expectedKey := lua.LNumber(1)
	for key, _ := table.Next(lua.LNil); key != lua.LNil; key, _ = table.Next(key) {
		if expectedKey != key {
			return false
		}
		expectedKey++
	}
	return true
}

func toYAML(L *lua.LState, visited map[*lua.LTable]bool, value lua.LValue) interface{} {
	switch value.Type() {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return lua.LVAsBool(value)
	case lua.LTNumber:
		num := float64(lua.LVAsNumber(value))
		intNum := int64(num)
		if num != float64(intNum) {
			return num
		}
		return intNum
	case lua.LTString:
		return lua.LVAsString(value)
	case lua.LTTable:
		valueTable := value.(*lua.LTable)
		if visited[valueTable] {
			L.RaiseError("nested table %s", valueTable)
		}
		visited[valueTable] = true
		if tableIsSlice(valueTable) {
			ret := make([]interface{}, 0, valueTable.Len())
			valueTable.ForEach(func(_ lua.LValue, tValue lua.LValue) {
				ret = append(ret, toYAML(L, visited, tValue))
			})
			return ret
		}
		ret := make(map[interface{}]interface{})
		valueTable.ForEach(func(tKey lua.LValue, tValue lua.LValue) {
			ret[toYAML(L, visited, tKey)] = toYAML(L, visited, tValue)
		})
		return ret
	default:
		L.RaiseError(fmt.Sprintf("cannot encode values with %s in them", value.Type()))
		return nil
	}
}

func fromYAML(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case int:
		return lua.LNumber(converted)
	case int64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(fromYAML(L, item))
		}
		return arr
	case map[interface{}]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(fromYAML(L, key), fromYAML(L, item))
		}
		return tbl
	case interface{}:
		if v, ok := converted.(bool); ok {
			return lua.LBool(v)
		}
		if v, ok := converted.(float64); ok {
			return lua.LNumber(v)
		}
		if v, ok := converted.(string); ok {
			return lua.LString(v)
		}
	}
	return lua.LNil
}
