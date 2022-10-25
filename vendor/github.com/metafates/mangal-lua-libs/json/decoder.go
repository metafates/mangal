package json

import (
	"encoding/json"
	"github.com/metafates/mangal-lua-libs/io"
	lua "github.com/yuin/gopher-lua"
)

const (
	jsonDecoderType = "json.Decoder"
)

func CheckJSONDecoder(L *lua.LState, n int) *json.Decoder {
	ud := L.CheckUserData(n)
	if decoder, ok := ud.Value.(*json.Decoder); ok {
		return decoder
	}
	L.ArgError(n, jsonDecoderType+" expected")
	return nil
}

func LVJSONDecoder(L *lua.LState, decoder *json.Decoder) lua.LValue {
	ud := L.NewUserData()
	ud.Value = decoder
	L.SetMetatable(ud, L.GetTypeMetatable(jsonDecoderType))
	return ud
}

func jsonDecoderDecode(L *lua.LState) int {
	decoder := CheckJSONDecoder(L, 1)
	L.Pop(L.GetTop())
	var value interface{}
	if err := decoder.Decode(&value); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(decode(L, value))
	return 1
}

func jsonDecoderInputOffset(L *lua.LState) int {
	decoder := CheckJSONDecoder(L, 1)
	L.Pop(L.GetTop())
	L.Push(lua.LNumber(decoder.InputOffset()))
	return 1
}

func jsonDecoderMore(L *lua.LState) int {
	decoder := CheckJSONDecoder(L, 1)
	L.Pop(L.GetTop())
	L.Push(lua.LBool(decoder.More()))
	return 1
}

func registerDecoder(L *lua.LState) {
	mt := L.NewTypeMetatable(jsonDecoderType)
	L.SetGlobal(jsonDecoderType, mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"decode":       jsonDecoderDecode,
		"input_offset": jsonDecoderInputOffset,
		"more":         jsonDecoderMore,
	}))
}

func registerJsonDecodedObject(L *lua.LState) {
	mt := L.NewTypeMetatable(jsonTableIsObject)
	mt.RawSetString(jsonTableIsObject, lua.LTrue)
}

func newJSONDecoder(L *lua.LState) int {
	reader := io.CheckIOReader(L, 1)
	L.Pop(L.GetTop())
	decoder := json.NewDecoder(reader)
	L.Push(LVJSONDecoder(L, decoder))
	return 1
}
