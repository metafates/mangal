package script

import (
	"context"
	"io"
	"reflect"

	"github.com/mangalorg/libmangal"
	"github.com/mangalorg/mangal/script/lib"
	"github.com/mangalorg/mangal/stringutil"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

type Variables = map[string]string

type Options struct {
	Client    *libmangal.Client
	Anilist   *libmangal.Anilist
	Variables Variables
}

func addVarsTable(state *lua.LState, variables Variables) {
	table := state.NewTable()
	for key, value := range variables {
		table.RawSetString(key, lua.LString(value))
	}

	state.SetGlobal("Vars", table)
}

func addLibraries(state *lua.LState, options lib.Options) {
	lib.Preload(state, options)
}

func Run(ctx context.Context, script io.Reader, options Options) error {
	state := lua.NewState()
	state.SetContext(ctx)
	config := luar.GetConfig(state)
	config.FieldNames = func(s reflect.Type, f reflect.StructField) []string {
		return []string{stringutil.CaseCamelToSnake(f.Name)}
	}
	config.MethodNames = func(t reflect.Type, m reflect.Method) []string {
		return []string{stringutil.CaseCamelToSnake(m.Name)}
	}

	addVarsTable(state, options.Variables)
	addLibraries(state, lib.Options{
		Client:  options.Client,
		Anilist: options.Anilist,
	})

	lFunction, err := state.Load(script, "script")
	if err != nil {
		return err
	}

	return state.CallByParam(lua.P{
		Fn:      lFunction,
		NRet:    1,
		Protect: true,
	})
}
