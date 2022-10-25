package json

import (
	"encoding/json"
	"github.com/metafates/mangal-lua-libs/io"
	lua "github.com/yuin/gopher-lua"
)

const (
	jsonEncoderType = "json.Encoder"
)

func CheckJSONEncoder(L *lua.LState, n int) *json.Encoder {
	ud := L.CheckUserData(n)
	if encoder, ok := ud.Value.(*json.Encoder); ok {
		return encoder
	}
	L.ArgError(n, jsonEncoderType+" expected")
	return nil
}

func LVJSONEncoder(L *lua.LState, encoder *json.Encoder) lua.LValue {
	ud := L.NewUserData()
	ud.Value = encoder
	L.SetMetatable(ud, L.GetTypeMetatable(jsonEncoderType))
	return ud
}

func jsonEncoderEncode(L *lua.LState) int {
	encoder := CheckJSONEncoder(L, 1)
	value := L.CheckAny(2)
	L.Pop(L.GetTop())
	err := encoder.Encode(jsonValue{
		LValue:  value,
		visited: make(map[*lua.LTable]bool),
	})
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

func jsonEncoderSetIndent(L *lua.LState) int {
	encoder := CheckJSONEncoder(L, 1)
	prefix := L.CheckString(2)
	indent := L.CheckString(3)
	L.Pop(L.GetTop())
	encoder.SetIndent(prefix, indent)
	return 0
}

func jsonEncoderSetEscapeHTML(L *lua.LState) int {
	encoder := CheckJSONEncoder(L, 1)
	on := L.CheckBool(2)
	L.Pop(L.GetTop())
	encoder.SetEscapeHTML(on)
	return 0
}

func registerJSONEncoder(L *lua.LState) {
	mt := L.NewTypeMetatable(jsonEncoderType)
	L.SetGlobal(jsonEncoderType, mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"encode":          jsonEncoderEncode,
		"set_indent":      jsonEncoderSetIndent,
		"set_escape_HTML": jsonEncoderSetEscapeHTML,
	}))
}

func newJSONEncoder(L *lua.LState) int {
	writer := io.CheckIOWriter(L, 1)
	L.Pop(L.GetTop())
	encoder := json.NewEncoder(writer)
	L.Push(LVJSONEncoder(L, encoder))
	return 1
}
