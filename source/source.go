package source

import (
	lua "github.com/yuin/gopher-lua"
)

type Source struct {
	Name  string
	state *lua.LState
}

func newSource(name string, state *lua.LState) (*Source, error) {
	return &Source{
		Name:  name,
		state: state,
	}, nil
}

func (s *Source) call(fn string, ret lua.LValueType, args ...lua.LValue) (lua.LValue, error) {
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

func (s *Source) Search(query string) ([]*Manga, error) {
	_, err := s.call(SearchMangaFn, lua.LTTable, lua.LString(query))

	if err != nil {
		return nil, err
	}

	table := s.state.CheckTable(-1)
	var i uint
	mangas := make([]*Manga, table.Len())

	table.ForEach(func(k lua.LValue, v lua.LValue) {
		if k.Type() != lua.LTNumber {
			s.state.RaiseError(SearchMangaFn + " was expected to return a table with numbers as keys, got " + k.Type().String() + " as a key")
		}

		if v.Type() != lua.LTTable {
			s.state.RaiseError(SearchMangaFn + " was expected to return a table with tables as values, got " + v.Type().String() + " as a value")
		}

		manga, err := mangaFromTable(v.(*lua.LTable))

		if err != nil {
			s.state.RaiseError(err.Error())
		}

		mangas[i] = manga
		i++
	})

	return mangas, nil
}

func (s *Source) Debug() error {
	_, err := s.call(DebugFn, lua.LTNil)
	return err
}
