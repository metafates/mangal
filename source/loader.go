package source

import (
	"errors"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/luamodules"
	"github.com/metafates/mangal/util"
	"github.com/spf13/afero"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

func IDfromName(name string) string {
	return name + " custom"
}

func LoadSource(path string, validate bool) (Source, error) {
	proto, err := Compile(path)
	if err != nil {
		return nil, err
	}

	state := lua.NewState()
	luamodules.PreloadAll(state)

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
				return nil, errors.New("required function " + fn + " is not defined in the source " + name)
			}
		}
	}

	source, err := newLuaSource(name, state)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func Compile(path string) (*lua.FunctionProto, error) {
	file, err := filesystem.Get().Open(path)

	if err != nil {
		return nil, err
	}

	defer func(file afero.File) {
		_ = file.Close()
	}(file)

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
