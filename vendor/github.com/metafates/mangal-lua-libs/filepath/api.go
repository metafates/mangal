// Package filepath implements golang filepath functionality for lua.
package filepath

import (
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

// Basename lua filepath.basename(path) returns the last element of path
func Basename(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.Base(path)))
	return 1
}

// Dir lua filepath.dir(path) returns all but the last element of path, typically the path's directory
func Dir(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.Dir(path)))
	return 1
}

// Ext lua filepath.ext(path) returns the file name extension used by path.
func Ext(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.Ext(path)))
	return 1
}

// Join lua fileapth.join(path, ...) joins any number of path elements into a single path, adding a Separator if necessary.
func Join(L *lua.LState) int {
	path := L.CheckString(1)
	for i := 2; i <= L.GetTop(); i++ {
		add := L.CheckAny(i).String()
		path = filepath.Join(path, add)
	}
	L.Push(lua.LString(path))
	return 1
}

//  Separator lua filepath.separator() OS-specific path separator
func Separator(L *lua.LState) int {
	L.Push(lua.LString(filepath.Separator))
	return 1
}

// ListSeparator lua filepath.list_separator() OS-specific path list separator
func ListSeparator(L *lua.LState) int {
	L.Push(lua.LString(filepath.ListSeparator))
	return 1
}

// filepath.glob(pattern) returns the names of all files matching pattern or nil if there is no matching file.
func Glob(L *lua.LState) int {
	pattern := L.CheckString(1)
	files, err := filepath.Glob(pattern)
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}
	result := L.CreateTable(len(files), 0)
	for _, file := range files {
		result.Append(lua.LString(file))
	}
	L.Push(result)
	return 1
}
