package lib

import lua "github.com/yuin/gopher-lua"

func LuaDoc() string {
	return Lib(lua.NewState(), Options{}).LuaDoc()
}
