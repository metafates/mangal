package xmlpath

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds xmlpath to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local xmlpath = require("xmlpath")
func Preload(L *lua.LState) {
	L.PreloadModule("xmlpath", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	node := L.NewTypeMetatable(`xmlpath_node_ud`)
	L.SetGlobal(`xmlpath_node_ud`, node)
	L.SetField(node, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"string": NodeToString,
	}))

	path := L.NewTypeMetatable(`xmlpath_path_ud`)
	L.SetGlobal(`xmlpath_path_ud`, path)
	L.SetField(path, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"iter": PathIter,
	}))

	iter := L.NewTypeMetatable(`xmlpath_iter_ud`)
	L.SetGlobal(`xmlpath_iter_ud`, iter)
	L.SetField(iter, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"node": IterNode,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"load":    Load,
	"compile": Compile,
}
