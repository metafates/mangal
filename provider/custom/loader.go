package custom

import (
	"fmt"
	libs "github.com/metafates/mangal-lua-libs"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

func IDfromName(name string) string {
	return name + " custom"
}

func LoadSource(path string, validate bool) (source.Source, error) {
	proto, err := Compile(path)
	if err != nil {
		return nil, err
	}

	state := lua.NewState()
	libs.Preload(state)

	lfunc := state.NewFunctionFromProto(proto)
	state.Push(lfunc)
	err = state.PCall(0, lua.MultRet, nil)
	if err != nil {
		return nil, err
	}

	name := util.FileStem(path)

	if validate {
		for _, fn := range mustHave {
			defined := state.GetGlobal(fn)

			if defined.Type() != lua.LTFunction {
				return nil, fmt.Errorf("required function %s is not defined in the luaSource %s", fn, name)
			}
		}
	}

	luaSource, err := newLuaSource(name, state)
	if err != nil {
		return nil, err
	}

	return luaSource, nil
}

func Compile(path string) (*lua.FunctionProto, error) {
	file, err := filesystem.Api().Open(path)

	if err != nil {
		return nil, err
	}

	defer util.Ignore(file.Close)

	chunk, err := parse.Parse(file, path)

	if err != nil {
		return nil, err
	}

	proto, err := lua.Compile(chunk, path)
	if err != nil {
		return nil, err
	}

	return proto, nil
}
