package yaml

import (
	"github.com/metafates/mangal-lua-libs/io"
	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v2"
)

const (
	yamlDecoderType = "yaml.Decoder"
)

func CheckYAMLDecoder(L *lua.LState, n int) *yaml.Decoder {
	ud := L.CheckUserData(n)
	if decoder, ok := ud.Value.(*yaml.Decoder); ok {
		return decoder
	}
	L.ArgError(n, yamlDecoderType+" expected")
	return nil
}

func LVYAMLDecoder(L *lua.LState, decoder *yaml.Decoder) lua.LValue {
	ud := L.NewUserData()
	ud.Value = decoder
	L.SetMetatable(ud, L.GetTypeMetatable(yamlDecoderType))
	return ud
}

func yamlDecoderDecode(L *lua.LState) int {
	decoder := CheckYAMLDecoder(L, 1)
	L.Pop(L.GetTop())
	var value interface{}
	if err := decoder.Decode(&value); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(fromYAML(L, value))
	return 1
}

func yamlDecoderSetStrict(L *lua.LState) int {
	decoder := CheckYAMLDecoder(L, 1)
	strict := L.CheckBool(2)
	L.Pop(L.GetTop())
	decoder.SetStrict(strict)
	return 0
}

func registerYAMLDecoder(L *lua.LState) {
	mt := L.NewTypeMetatable(yamlDecoderType)
	L.SetGlobal(yamlDecoderType, mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"decode":     yamlDecoderDecode,
		"set_strict": yamlDecoderSetStrict,
	}))
}

func newYAMLDecoder(L *lua.LState) int {
	reader := io.CheckIOReader(L, 1)
	L.Pop(L.GetTop())
	decoder := yaml.NewDecoder(reader)
	L.Push(LVYAMLDecoder(L, decoder))
	return 1
}
