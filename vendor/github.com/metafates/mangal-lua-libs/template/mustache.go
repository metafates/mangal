package template

import (
	"sync"

	mustache "github.com/cbroglie/mustache"
	gluamapper "github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type luaMustache struct {
	sync.Mutex
	mapper *gluamapper.Mapper
}

func init() {
	nameFunc := func(name string) string { return name }
	RegisterTemplateEngine(`mustache`, &luaMustache{
		mapper: gluamapper.NewMapper(gluamapper.Option{
			NameFunc: nameFunc,
		}),
	})
}

func (t *luaMustache) Render(data string, context *lua.LTable) (string, error) {
	var values map[string]interface{}
	if err := t.mapper.Map(context, &values); err != nil {
		return "", err
	}
	return mustache.Render(data, values)
}
