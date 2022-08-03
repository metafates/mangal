package source

import (
	lua "github.com/yuin/gopher-lua"
	"strconv"
)

type LuaSource struct {
	name  string
	state *lua.LState
}

func (s *LuaSource) Name() string {
	return s.name
}

func newLuaSource(name string, state *lua.LState) (*LuaSource, error) {
	return &LuaSource{
		name:  name,
		state: state,
	}, nil
}

func (s *LuaSource) call(fn string, ret lua.LValueType, args ...lua.LValue) (lua.LValue, error) {
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

func (s *LuaSource) Search(query string) ([]*Manga, error) {
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

		index, err := strconv.ParseUint(k.String(), 10, 16)
		if err != nil {
			s.state.RaiseError(SearchMangaFn + " was expected to return a table with unsigned integers as keys. " + err.Error())
		}

		manga, err := mangaFromTable(v.(*lua.LTable), uint16(index))

		if err != nil {
			s.state.RaiseError(err.Error())
		}

		mangas[i] = manga
		i++
	})

	return mangas, nil
}

func (s *LuaSource) ChaptersOf(manga *Manga) ([]*Chapter, error) {
	_, err := s.call(MangaChaptersFn, lua.LTTable, lua.LString(manga.URL))

	if err != nil {
		return nil, err
	}

	table := s.state.CheckTable(-1)
	var i uint
	chapters := make([]*Chapter, table.Len())

	table.ForEach(func(k lua.LValue, v lua.LValue) {
		if k.Type() != lua.LTNumber {
			s.state.RaiseError(MangaChaptersFn + " was expected to return a table with numbers as keys, got " + k.Type().String() + " as a key")
		}

		if v.Type() != lua.LTTable {
			s.state.RaiseError(MangaChaptersFn + " was expected to return a table with tables as values, got " + v.Type().String() + " as a value")
		}

		index, err := strconv.ParseUint(k.String(), 10, 16)
		if err != nil {
			s.state.RaiseError(MangaChaptersFn + " was expected to return a table with unsigned integers as keys. " + err.Error())
		}

		chapter, err := chapterFromTable(v.(*lua.LTable), manga, uint16(index))

		if err != nil {
			s.state.RaiseError(err.Error())
		}

		chapters[i] = chapter
		i++
	})

	return chapters, nil
}

func (s *LuaSource) PagesOf(chapter *Chapter) ([]*Page, error) {
	_, err := s.call(ChapterPagesFn, lua.LTTable, lua.LString(chapter.URL))

	if err != nil {
		return nil, err
	}

	table := s.state.CheckTable(-1)
	var i uint
	pages := make([]*Page, table.Len())

	table.ForEach(func(k lua.LValue, v lua.LValue) {
		if k.Type() != lua.LTNumber {
			s.state.RaiseError(ChapterPagesFn + " was expected to return a table with numbers as keys, got " + k.Type().String() + " as a key")
		}

		if v.Type() != lua.LTTable {
			s.state.RaiseError(ChapterPagesFn + " was expected to return a table with tables as values, got " + v.Type().String() + " as a value")
		}

		page, err := pageFromTable(v.(*lua.LTable), chapter)

		if err != nil {
			s.state.RaiseError(err.Error())
		}

		pages[i] = page
		i++
	})

	return pages, nil
}

func (s *LuaSource) Debug() error {
	_, err := s.call(DebugFn, lua.LTNil)
	return err
}
