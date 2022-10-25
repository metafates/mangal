// Port dustin/go-humanize for gopher-lua
package humanize

import (
	"time"

	humanize "github.com/dustin/go-humanize"
	lua "github.com/yuin/gopher-lua"
)

// Time lua humanize.time(number) return string
func Time(L *lua.LState) int {
	then := time.Unix(L.CheckInt64(1), 0)
	L.Push(lua.LString(humanize.Time(then)))
	return 1
}

// IBytes lua humanize.ibytes(number) return string
func IBytes(L *lua.LState) int {
	bytes := L.CheckInt64(1)
	if bytes < 0 {
		L.ArgError(1, "must be positive")
	}
	L.Push(lua.LString(humanize.IBytes(uint64(bytes))))
	return 1
}

// ParseBytes lua humanize.parse_bytes(string) returns (number, err)
func ParseBytes(L *lua.LState) int {
	data := L.CheckString(1)
	size, err := humanize.ParseBytes(data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(size))
	return 1
}

// SI lua humanize.si(number, string) return string
func SI(L *lua.LState) int {
	value := L.CheckNumber(1)
	input := float64(value)
	unit := L.CheckString(2)
	L.Push(lua.LString(humanize.SI(input, unit)))
	return 1
}
