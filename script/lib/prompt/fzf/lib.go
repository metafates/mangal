package fzf

import (
	"github.com/ktr0731/go-fuzzyfinder"
	luadoc "github.com/mangalorg/gopher-luadoc"
	"github.com/mangalorg/mangal/script/lib/util"
	lua "github.com/yuin/gopher-lua"
)

const libName = "fzf"

func Lib() *luadoc.Lib {
	const generic = "T"

	return &luadoc.Lib{
		Name:        libName,
		Description: "fuzzy-finding with an fzf-like terminal user interface.",
		Funcs: []*luadoc.Func{
			{
				Name:        "find_one",
				Description: "Find displays a UI that provides fuzzy finding against the provided list",
				Value:       findOne,
				Generics:    []string{generic},
				Params: []*luadoc.Param{
					{
						Name:        "items",
						Description: "",
						Type:        luadoc.List(generic),
					},
					{
						Name:        "to_string",
						Description: "Function that is used to convert item to the string representation",
						Type: luadoc.Func{
							Name:        libName,
							Description: "",
							Params: []*luadoc.Param{
								{
									Name: "item",
									Type: generic,
								},
							},
							Returns: []*luadoc.Param{
								{
									Name: "representation",
									Type: luadoc.String,
								},
							},
						}.AsType(),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "item",
						Description: "",
						Type:        generic,
					},
				},
			},
		},
	}
}

func findOne(state *lua.LState) int {
	var values []lua.LValue

	state.CheckTable(1).ForEach(func(_, value lua.LValue) {
		values = append(values, value)
	})

	toString := state.CheckFunction(2)

	index, err := fuzzyfinder.Find(values, func(i int) string {
		state.Push(toString)
		state.Push(values[i])

		err := state.PCall(1, 1, nil)
		util.Must(state, err)

		return state.Get(-1).String()
	})

	util.Must(state, err)

	state.Push(values[index])
	return 1
}
