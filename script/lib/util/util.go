package util

import (
	"fmt"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	lua "github.com/yuin/gopher-lua"
)

func NewUserData[T any](state *lua.LState, t T, typeName string) *lua.LUserData {
	ud := state.NewUserData()
	ud.Value = t

	state.SetMetatable(ud, state.GetTypeMetatable(typeName))

	return ud
}

func Push[T any](state *lua.LState, t T, typeName string) {
	state.Push(NewUserData(state, t, typeName))
}

func Check[T any](state *lua.LState, n int) T {
	t, ok := state.CheckUserData(n).Value.(T)
	if !ok {
		state.ArgError(n, fmt.Sprintf("%T expected", t))
	}

	return t
}

func Must(state *lua.LState, err error) {
	if err != nil {
		state.RaiseError("%s", err)
	}
}

func SliceToTable[T any](state *lua.LState, items []T, convert func(T) lua.LValue) *lua.LTable {
	table := state.NewTable()

	for _, item := range items {
		table.Append(convert(item))
	}

	return table
}

func ToGoValue(value lua.LValue) any {
	switch value.Type() {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return bool(value.(lua.LBool))
	case lua.LTNumber:
		return float64(value.(lua.LNumber))
	case lua.LTString:
		return string(value.(lua.LString))
	case lua.LTTable:
		table := value.(*lua.LTable)
		om := orderedmap.New[any, any]()

		var asMap = make(map[any]any)

		table.ForEach(func(key lua.LValue, value lua.LValue) {
			k, v := ToGoValue(key), ToGoValue(value)
			asMap[k] = v
			om.Set(k, v)
		})

		// check if we can convert table to slice.
		// if not, return as map.
		var (
			prev    float64 = 0
			asSlice []any
		)
		for pair := om.Oldest(); pair != nil; pair = pair.Next() {
			asNum, ok := pair.Key.(float64)
			if !ok || asNum != prev+1 {
				return asMap
			}

			prev = asNum
			asSlice = append(asSlice, pair.Value)
		}

		return asSlice
	case lua.LTUserData:
		return value.(*lua.LUserData).Value
	default:
		return nil
	}
}
