// Package goos implements golang package os functionality for lua.
package goos

import (
	"os"

	lua "github.com/yuin/gopher-lua"
)

// Stat lua goos.stat(filename) returns (table, err)
func Stat(L *lua.LState) int {
	filename := L.CheckString(1)
	stat, err := os.Stat(filename)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.NewTable()
	result.RawSetString(`is_dir`, lua.LBool(stat.IsDir()))
	result.RawSetString(`size`, lua.LNumber(stat.Size()))
	result.RawSetString(`mod_time`, lua.LNumber(stat.ModTime().Unix()))
	result.RawSetString(`mode`, lua.LString(stat.Mode().String()))
	L.Push(result)
	return 1
}

// Hostname lua goos.hostname() returns (string, error)
func Hostname(L *lua.LState) int {
	hostname, err := os.Hostname()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(hostname))
	return 1
}

// Getpagesize lua goos.pagesize() return number
func Getpagesize(L *lua.LState) int {
	L.Push(lua.LNumber(os.Getpagesize()))
	return 1
}

// MkdirAll lua goos.mkdir_all() return err
func MkdirAll(L *lua.LState) int {
	err := os.MkdirAll(L.CheckString(1), 0755)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
