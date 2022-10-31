package custom

import (
	"fmt"
	"github.com/metafates/mangal/source"
	lua "github.com/yuin/gopher-lua"
)

type luaSource struct {
	name  string
	state *lua.LState
	cache struct {
		mangas   *cacher[[]*source.Manga]
		chapters *cacher[[]*source.Chapter]
	}
}

func (s *luaSource) Name() string {
	return s.name
}

func newLuaSource(name string, state *lua.LState) (*luaSource, error) {
	s := &luaSource{
		name:  name,
		state: state,
	}

	cacheName := func(cacheFor string) string {
		return fmt.Sprintf("%s_%s", s.ID(), cacheFor)
	}

	s.cache.mangas = newCacher[[]*source.Manga](cacheName("mangas"))
	s.cache.chapters = newCacher[[]*source.Chapter](cacheName("chapters"))

	return s, nil
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
