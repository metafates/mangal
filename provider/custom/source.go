package custom

import (
	"github.com/metafates/mangal/source"
	lua "github.com/yuin/gopher-lua"
)

type luaSource struct {
	name           string
	state          *lua.LState
	cachedMangas   map[string][]*source.Manga
	cachedChapters map[string][]*source.Chapter
	cachedPages    map[string][]*source.Page
}

func (s *luaSource) Name() string {
	return s.name
}

func newLuaSource(name string, state *lua.LState) (*luaSource, error) {
	return &luaSource{
		name:           name,
		state:          state,
		cachedMangas:   make(map[string][]*source.Manga),
		cachedChapters: make(map[string][]*source.Chapter),
		cachedPages:    make(map[string][]*source.Page),
	}, nil
}

func (s *luaSource) call(fn string, ret lua.LValueType, args ...lua.LValue) (lua.LValue, error) {
	err := s.state.CallByParam(lua.P{
		Fn:      s.state.GetGlobal(fn),
		NRet:    1,
		Protect: true,
	}, args...)

	if err != nil {
		return nil, err
	}

	val := s.state.Get(-1)

	if val.Type() != ret {
		s.state.RaiseError(fn + " was expected to return a " + ret.String() + ", got " + val.Type().String())
	}

	return val, nil
}

func (s *luaSource) ID() string {
	return IDfromName(s.name)
}
