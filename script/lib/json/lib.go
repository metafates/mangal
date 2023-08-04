package json

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/itchyny/gojq"
	json "github.com/json-iterator/go"
	luadoc "github.com/mangalorg/gopher-luadoc"
	"github.com/mangalorg/mangal/script/lib/util"
	lua "github.com/yuin/gopher-lua"
)

const libName = "json"

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        libName,
		Description: "JSON related functionality",
		Funcs: []*luadoc.Func{
			{
				Name:        "print",
				Description: "",
				Value:       jsonPrint,
				Params: []*luadoc.Param{
					{
						Name:        "value",
						Description: "",
						Type:        luadoc.Any,
					},
					{
						Name:        "pattern",
						Description: "JQ pattern to filter out JSON",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
		},
	}
}

func marshal(value any) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func jsonPrint(state *lua.LState) int {
	value := util.ToGoValue(state.CheckAny(1))

	marshalled, err := marshal(value)
	util.Must(state, err)

	JQPattern := state.OptString(2, "")

	if JQPattern == "" {
		fmt.Println(marshalled)
		return 0
	}

	query, err := gojq.Parse(JQPattern)
	util.Must(state, err)

	deadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	var iter gojq.Iter

	if reflect.ValueOf(value).Kind() == reflect.Slice {
		var v []any
		if err := json.Unmarshal([]byte(marshalled), &v); err != nil {
			state.RaiseError(err.Error())
			return 0
		}
		util.Must(state, err)

		iter = query.RunWithContext(ctx, v)
	} else {
		var v map[string]any
		err := json.Unmarshal([]byte(marshalled), &v)
		util.Must(state, err)

		iter = query.RunWithContext(ctx, v)
	}

	for {
		next, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := next.(error); ok {
			state.RaiseError(err.Error())
			return 0
		}

		toPrint, err := marshal(next)
		util.Must(state, err)

		fmt.Println(toPrint)
	}

	return 0
}
