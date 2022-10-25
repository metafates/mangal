// Package runtime implements golang runtime functionality for lua.
package runtime

import (
	"runtime"

	lua "github.com/yuin/gopher-lua"
)

// GOOS lua runtime.goos() return string
func GOOS(L *lua.LState) int {
	L.Push(lua.LString(runtime.GOOS))
	return 1
}

// GOARCH lua runtime.goarch() return string
func GOARCH(L *lua.LState) int {
	L.Push(lua.LString(runtime.GOARCH))
	return 1
}
