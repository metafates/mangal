// Package strings implements golang package montanaflynn/stats functionality for lua.
package stats

import (
	"fmt"

	gostats "github.com/montanaflynn/stats"
	lua "github.com/yuin/gopher-lua"
)

// get float slice from table
func getFloatSliceFromTable(L *lua.LState, n int) ([]float64, error) {
	tbl := L.CheckTable(n)
	data := make([]float64, tbl.Len())
	var err error
	tbl.ForEach(func(k lua.LValue, v lua.LValue) {
		value, ok := v.(lua.LNumber)
		if !ok {
			err = fmt.Errorf("only table of numbers is supported")
			return
		}
		data = append(data, float64(value))
	})
	return data, err
}

// Median lua stats.median(table): port of go montanaflynn/stats.Median() returns value and error
func Median(L *lua.LState) int {
	data, err := getFloatSliceFromTable(L, 1)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result, err := gostats.Median(data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(result))
	return 1
}

// Percentile lua stats.median(table, percentile): port of go montanaflynn/stats.Percentile() returns value and error
func Percentile(L *lua.LState) int {
	data, err := getFloatSliceFromTable(L, 1)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	percentile := L.CheckNumber(2)
	result, err := gostats.Percentile(data, float64(percentile))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(result))
	return 1
}

// StandardDeviation lua stats.median(table, percentile): port of go montanaflynn/stats.StandardDeviation() returns value and error
func StandardDeviation(L *lua.LState) int {
	data, err := getFloatSliceFromTable(L, 1)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result, err := gostats.StandardDeviation(data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(result))
	return 1
}
