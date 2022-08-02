package source

import (
	"bufio"
	"errors"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/luamodules"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"os"
	"path/filepath"
)

const sourceExtension = ".lua"

func LoadSource(path string) (*Source, error) {
	proto, err := compileSource(path)
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

	base := filepath.Base(path)
	name := base[:len(base)-len(sourceExtension)]

	for _, fn := range mustHave {
		defined := state.GetGlobal(fn)

		if defined.Type() != lua.LTFunction {
			return nil, errors.New("required function " + fn + " is not defined in the source " + name)
		}
	}

	source, err := newSource(name, state)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func AvailableSources() ([]string, error) {
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

func compileSource(path string) (*lua.FunctionProto, error) {
	if exists, err := filesystem.Get().Exists(path); !exists || err != nil {
		return nil, errors.New("source file does not exist")
	}

	file, err := filesystem.Get().Open(path)

	if err != nil {
		return nil, err
	}

	defer func(file afero.File) {
		_ = file.Close()
	}(file)

	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, path)

	if err != nil {
		return nil, err
	}

	proto, err := lua.Compile(chunk, path)
	if err != nil {
		return nil, err
	}

	return proto, nil
}
