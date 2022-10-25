package strings

import (
	lio "github.com/metafates/mangal-lua-libs/io"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

const (
	stringsBuilderType = "strings.Builder"
)

func CheckStringsBuilder(L *lua.LState, n int) *strings.Builder {
	ud := L.CheckUserData(n)
	if builder, ok := ud.Value.(*strings.Builder); ok {
		return builder
	}
	L.ArgError(n, stringsBuilderType+" expected")
	return nil
}

func LVStringsBuilder(L *lua.LState, builder *strings.Builder) lua.LValue {
	ud := L.NewUserData()
	ud.Value = builder
	L.SetMetatable(ud, L.GetTypeMetatable(stringsBuilderType))
	return ud
}

func stringsBuilderString(L *lua.LState) int {
	builder := CheckStringsBuilder(L, 1)
	s := builder.String()
	L.Push(lua.LString(s))
	return 1
}

func newStringsBuilder(L *lua.LState) int {
	builder := &strings.Builder{}
	L.Push(LVStringsBuilder(L, builder))
	return 1
}

func registerStringsBuilder(L *lua.LState) {
	mt := L.NewTypeMetatable(stringsBuilderType)
	L.SetGlobal(stringsBuilderType, mt)
	writerTable := lio.WriterFuncTable(L)
	L.SetField(writerTable, "string", L.NewFunction(stringsBuilderString))
	L.SetField(mt, "__index", writerTable)
}
