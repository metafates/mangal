package source

import (
	"bufio"
	"errors"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/luamodules"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"io"
	"os"
	"path/filepath"
)

const sourceExtension = ".lua"

func LoadSource(name string, proto *lua.FunctionProto) (Source, error) {
	state := lua.NewState()
	luamodules.PreloadAll(state)

	lfunc := state.NewFunctionFromProto(proto)
	state.Push(lfunc)
	err := state.PCall(0, lua.MultRet, nil)
	if err != nil {
		return nil, err
	}

	for _, fn := range mustHave {
		defined := state.GetGlobal(fn)

		if defined.Type() != lua.LTFunction {
			return nil, errors.New("required function " + fn + " is not defined in the source " + name)
		}
	}

	source, err := newLuaSource(name, state)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func AvailableCustomSources() ([]string, error) {
	if exists := lo.Must(filesystem.Get().Exists(viper.GetString(config.SourcesPath))); !exists {
		return nil, errors.New("sources directory does not exist")
	}

	files, err := filesystem.Get().ReadDir(viper.GetString(config.SourcesPath))

	if err != nil {
		return nil, err
	}

	return lo.FilterMap(files, func(f os.FileInfo, _ int) (string, bool) {
		if filepath.Ext(f.Name()) == sourceExtension {
			return filepath.Join(viper.GetString(config.SourcesPath), f.Name()), true
		}

		return "", false
	}), nil
}

func Compile(name string, script io.Reader) (*lua.FunctionProto, error) {
	chunk, err := parse.Parse(bufio.NewReader(script), name)

	if err != nil {
		return nil, err
	}

	proto, err := lua.Compile(chunk, name)
	if err != nil {
		return nil, err
	}

	return proto, nil
}
