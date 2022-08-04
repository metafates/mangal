package source

import (
	"errors"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/luamodules"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"os"
	"path/filepath"
	"strings"
)

const sourceExtension = ".lua"

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

	if validate {
		name := strings.TrimSuffix(filepath.Base(path), sourceExtension)
		for _, fn := range mustHave {
			defined := state.GetGlobal(fn)

			if defined.Type() != lua.LTFunction {
				return nil, errors.New("required function " + fn + " is not defined in the source " + name)
			}
		}
	}

	source, err := newLuaSource(path, state)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func AvailableCustomSources() (map[string]string, error) {
	if exists := lo.Must(filesystem.Get().Exists(viper.GetString(config.SourcesPath))); !exists {
		return nil, errors.New("sources directory does not exist")
	}

	files, err := filesystem.Get().ReadDir(viper.GetString(config.SourcesPath))

	if err != nil {
		return nil, err
	}

	sources := make(map[string]string)
	paths := lo.FilterMap(files, func(f os.FileInfo, _ int) (string, bool) {
		if filepath.Ext(f.Name()) == sourceExtension {
			return filepath.Join(viper.GetString(config.SourcesPath), f.Name()), true
		}
		return "", false
	})

	for _, path := range paths {
		name := strings.TrimSuffix(filepath.Base(path), sourceExtension) + " " + icon.Lua()
		sources[name] = path
	}

	return sources, nil
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
