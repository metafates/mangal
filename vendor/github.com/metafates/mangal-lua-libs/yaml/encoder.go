package yaml

import (
	"github.com/metafates/mangal-lua-libs/io"
	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v2"
)

const (
	yamlEncoderType = "yaml.Encoder"
)

func CheckYAMLEncoder(L *lua.LState, n int) *yaml.Encoder {
	ud := L.CheckUserData(n)
	if encoder, ok := ud.Value.(*yaml.Encoder); ok {
		return encoder
	}
	L.ArgError(n, yamlEncoderType+" expected")
	return nil
}

func LVYAMLEncoder(L *lua.LState, encoder *yaml.Encoder) lua.LValue {
	ud := L.NewUserData()
	ud.Value = encoder
	L.SetMetatable(ud, L.GetTypeMetatable(yamlEncoderType))
	return ud
}

func yamlEncoderEncode(L *lua.LState) int {
	encoder := CheckYAMLEncoder(L, 1)
	arg := L.CheckAny(2)
	L.Pop(L.GetTop())
	var value interface{}
	err := L.GPCall(func(L *lua.LState) int {
		visited := make(map[*lua.LTable]bool)
		value = toYAML(L, visited, arg)
		return 0
	}, lua.LNil)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	if err = encoder.Encode(value); err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

func registerYAMLEncoder(L *lua.LState) {
	mt := L.NewTypeMetatable(yamlEncoderType)
	L.SetGlobal(yamlEncoderType, mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"encode": yamlEncoderEncode,
	}))
}

func newYAMLEncoder(L *lua.LState) int {
	writer := io.CheckIOWriter(L, 1)
	L.Pop(L.GetTop())
	encoder := yaml.NewEncoder(writer)
	L.Push(LVYAMLEncoder(L, encoder))
	return 1
}
