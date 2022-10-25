// Package ioutil implements golang package ioutil functionality for lua.
package ioutil

import (
	lio "github.com/metafates/mangal-lua-libs/io"
	"io"
	"io/ioutil"

	lua "github.com/yuin/gopher-lua"
)

// ReadFile lua ioutil.read_file(filepath) reads the file named by filename and returns the contents, returns (string,error)
func ReadFile(L *lua.LState) int {
	filename := L.CheckString(1)
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		L.Push(lua.LString(data))
		return 1
	} else {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
}

// WriteFile lua ioutil.write_file(filepath, data) reads the file named by filename and returns the contents, returns (string,error)
func WriteFile(L *lua.LState) int {
	filename := L.CheckString(1)
	data := L.CheckString(2)
	err := ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

func Copy(L *lua.LState) int {
	writer := lio.CheckIOWriter(L, 1)
	reader := lio.CheckIOReader(L, 2)
	L.Pop(L.GetTop())
	if _, err := io.Copy(writer, reader); err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

func CopyN(L *lua.LState) int {
	writer := lio.CheckIOWriter(L, 1)
	reader := lio.CheckIOReader(L, 2)
	n := L.CheckInt64(3)
	L.Pop(L.GetTop())
	if _, err := io.CopyN(writer, reader, n); err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
