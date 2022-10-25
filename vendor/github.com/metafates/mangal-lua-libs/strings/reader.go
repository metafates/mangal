package strings

import (
	"github.com/metafates/mangal-lua-libs/io"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

const (
	stringsReaderType = "strings.Reader"
)

func CheckStringsReader(L *lua.LState, n int) *strings.Reader {
	ud := L.CheckUserData(n)
	if reader, ok := ud.Value.(*strings.Reader); ok {
		return reader
	}
	L.ArgError(n, stringsReaderType+" expected")
	return nil
}

func LVStringsReader(L *lua.LState, reader *strings.Reader) lua.LValue {
	ud := L.NewUserData()
	ud.Value = reader
	L.SetMetatable(ud, L.GetTypeMetatable(stringsReaderType))
	return ud
}

func newStringsReader(L *lua.LState) int {
	s := L.CheckString(1)
	L.Pop(L.GetTop())
	reader := strings.NewReader(s)
	L.Push(LVStringsReader(L, reader))
	return 1
}

func registerStringsReader(L *lua.LState) {
	mt := L.NewTypeMetatable(stringsReaderType)
	L.SetGlobal(stringsReaderType, mt)
	readerTable := io.ReaderFuncTable(L)
	L.SetField(mt, "__index", readerTable)
}
