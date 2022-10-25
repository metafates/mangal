package io

import (
	"errors"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"io"
	"io/ioutil"
)

type luaIOWrapper struct {
	ls  *lua.LState
	obj lua.LValue

	readMethod  *lua.LFunction
	writeMethod *lua.LFunction
	seekMethod  *lua.LFunction
	closeMethod *lua.LFunction
}

//NewLuaIOWrapper creates a new luaIOWrapper atop the lua io object
func NewLuaIOWrapper(L *lua.LState, io lua.LValue) *luaIOWrapper {
	ret := &luaIOWrapper{
		ls:  L,
		obj: io,
	}
	ret.readMethod, _ = L.GetField(io, "read").(*lua.LFunction)
	ret.writeMethod, _ = L.GetField(io, "write").(*lua.LFunction)
	ret.seekMethod, _ = L.GetField(io, "seek").(*lua.LFunction)
	ret.closeMethod, _ = L.GetField(io, "close").(*lua.LFunction)
	return ret
}

//CheckIOWriter tries to cast to UserData and to io.Writer, otherwise it wraps and checks for "write" method
func CheckIOWriter(L *lua.LState, n int) io.Writer {
	any := L.CheckAny(n)
	if ud, ok := any.(*lua.LUserData); ok {
		if writer, ok := ud.Value.(io.Writer); ok {
			return writer
		}
	}
	wrapped := NewLuaIOWrapper(L, any)
	if wrapped.writeMethod == nil {
		L.ArgError(n, "expected writer")
		return nil
	}
	return wrapped
}

//CheckIOReader tries to cast to UserData and to io.Reader, otherwise it wraps and checks for "read" method
func CheckIOReader(L *lua.LState, n int) io.Reader {
	any := L.CheckAny(n)
	if ud, ok := any.(*lua.LUserData); ok {
		if reader, ok := ud.Value.(io.Reader); ok {
			return reader
		}
	}
	wrapped := NewLuaIOWrapper(L, any)
	if wrapped.readMethod == nil {
		L.ArgError(n, "expected reader")
		return nil
	}
	return wrapped
}

func (l *luaIOWrapper) Read(p []byte) (n int, err error) {
	if l.readMethod == nil {
		return 0, errors.New("object does not have read method")
	}
	n = len(p)

	L := l.ls
	L.Push(l.readMethod)
	L.Push(l.obj)
	L.Push(lua.LNumber(n))
	if err = L.PCall(2, 1, nil); err != nil {
		n = 0
		return
	}
	result := L.Get(1)
	L.Pop(L.GetTop())
	if result.Type() == lua.LTNil {
		return 0, io.EOF
	}
	readString := lua.LVAsString(result)
	data := []byte(readString)
	n = copy(p, data)
	return
}

func (l *luaIOWrapper) Write(p []byte) (n int, err error) {
	if l.writeMethod == nil {
		return 0, errors.New("object does not have write method")
	}
	n = len(p)
	L := l.ls
	L.Push(l.writeMethod)
	L.Push(l.obj)
	L.Push(lua.LString(p))
	err = L.PCall(2, 0, nil)
	return
}

func (l *luaIOWrapper) Seek(offset int64, whence int) (int64, error) {
	if l.seekMethod == nil {
		return 0, errors.New("object does not have seek method")
	}
	var luaWhence string
	switch whence {
	case io.SeekStart:
		luaWhence = "set"
	case io.SeekEnd:
		luaWhence = "end"
	case io.SeekCurrent:
		luaWhence = "cur"
	default:
		return 0, fmt.Errorf("unknown whence: %d", whence)
	}

	L := l.ls
	L.Push(l.seekMethod)
	L.Push(l.obj)
	L.Push(lua.LString(luaWhence))
	L.Push(lua.LNumber(offset))
	if err := L.PCall(3, 1, nil); err != nil {
		return 0, err
	}
	ret := L.CheckNumber(1)
	L.Pop(L.GetTop())
	return int64(ret), nil
}

func (l *luaIOWrapper) Close() error {
	if l.closeMethod == nil {
		return errors.New("object does not have close method")
	}
	L := l.ls
	L.Push(l.closeMethod)
	L.Push(l.obj)
	return L.PCall(1, 0, nil)
}

func IOWriterWrite(L *lua.LState) int {
	writer := CheckIOWriter(L, 1)
	var toWrite []string
	for i := 2; i <= L.GetTop(); i++ {
		toWrite = append(toWrite, L.CheckString(i))
	}
	L.Pop(L.GetTop())
	for i, s := range toWrite {
		if _, err := io.WriteString(writer, s); err != nil {
			L.ArgError(i+2, err.Error())
			return 0
		}
	}
	return 0
}

func IOWriterClose(L *lua.LState) int {
	writer := CheckIOWriter(L, 1)
	L.Pop(L.GetTop())
	if closer, ok := writer.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			L.RaiseError("%v", err)
		}
	}
	return 0
}

func IOReaderRead(L *lua.LState) int {
	reader := CheckIOReader(L, 1)
	lFormat := L.Get(2)
	if lFormat == lua.LNil {
		lFormat = lua.LString("*l")
	}
	L.Pop(L.GetTop())

	if lFormat.Type() == lua.LTNumber {
		num := int(lua.LVAsNumber(lFormat))
		if num <= 0 {
			L.Push(lua.LString(""))
			return 1
		}
		buf := make([]byte, num)
		numRead, err := reader.Read(buf)
		if err == io.EOF {
			L.Push(lua.LNil)
			return 1
		}
		if err != nil {
			L.RaiseError("%v", err)
			return 0
		}
		if numRead < num {
			buf = buf[:numRead]
		}
		L.Push(lua.LString(buf))
		return 1
	}

	format := lua.LVAsString(lFormat)
	if len(format) >= 2 && format[0] == '*' {
		switch format[1] {
		case 'n':
			var num lua.LNumber
			_, err := fmt.Fscan(reader, &num)
			if err == io.EOF {
				L.Push(lua.LNumber(0))
				return 1
			}
			if err != nil {
				L.RaiseError("%v", err)
				return 0
			}
			L.Push(num)
			return 1
		case 'a':
			data, err := ioutil.ReadAll(reader)
			if err == io.EOF {
				L.Push(lua.LString(""))
				return 1
			}
			if err != nil {
				L.RaiseError("%v", err)
				return 0
			}
			L.Push(lua.LString(data))
			return 1
		case 'l':
			var line lua.LString
			_, err := fmt.Fscanln(reader, &line)
			if err == io.EOF {
				L.Push(lua.LNil)
				return 1
			}
			if err != nil {
				L.RaiseError("%v", err)
				return 0
			}
			L.Push(line)
			return 1
		}
	}

	L.ArgError(2, "unknown fmt string")
	return 0
}

func IOReaderClose(L *lua.LState) int {
	reader := CheckIOReader(L, 1)
	L.Pop(L.GetTop())
	if closer, ok := reader.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			L.RaiseError("%v", err)
		}
	}
	return 0
}

func WriterFuncTable(L *lua.LState) *lua.LTable {
	table := L.NewTable()
	L.SetFuncs(table, map[string]lua.LGFunction{
		"write": IOWriterWrite,
		"close": IOWriterClose,
	})
	return table
}

func ReaderFuncTable(L *lua.LState) *lua.LTable {
	table := L.NewTable()
	L.SetFuncs(table, map[string]lua.LGFunction{
		"read":  IOReaderRead,
		"close": IOReaderClose,
	})
	return table
}
