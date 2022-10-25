package template

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds template to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local template = require("template")
func Preload(L *lua.LState) {
	L.PreloadModule("template", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	template_ud := L.NewTypeMetatable(`template_ud`)
	L.SetGlobal(`template_ud`, template_ud)
	L.SetField(template_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"render":      Render,
		"render_file": RenderFile,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"choose": Choose,
}
