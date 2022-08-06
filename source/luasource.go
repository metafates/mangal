package source

import (
	lua "github.com/yuin/gopher-lua"
	"strconv"
)

type LuaSource struct {
	name           string
	state          *lua.LState
	cachedMangas   map[string][]*Manga
	cachedChapters map[string][]*Chapter
	cachedPages    map[string][]*Page
}

func (s *LuaSource) Name() string {
	return s.name
}

func newLuaSource(name string, state *lua.LState) (*LuaSource, error) {
	return &LuaSource{
		name:           name,
		state:          state,
		cachedMangas:   make(map[string][]*Manga),
		cachedChapters: make(map[string][]*Chapter),
		cachedPages:    make(map[string][]*Page),
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
	if cached, ok := s.cachedMangas[query]; ok {
		return cached, nil
	}

	_, err := s.call(SearchMangaFn, lua.LTTable, lua.LString(query))

	if err != nil {
		return nil, err
	}

	table := s.state.CheckTable(-1)
	mangas := make([]*Manga, 0)

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

		manga.SourceID = s.ID()

		mangas = append(mangas, manga)
	})

	s.cachedMangas[query] = mangas
	return mangas, nil
}

func (s *LuaSource) ChaptersOf(manga *Manga) ([]*Chapter, error) {
	if cached, ok := s.cachedChapters[manga.URL]; ok {
		return cached, nil
	}

	_, err := s.call(MangaChaptersFn, lua.LTTable, lua.LString(manga.URL))

	if err != nil {
		return nil, err
	}

	table := s.state.CheckTable(-1)
	chapters := make([]*Chapter, 0)

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

		chapter.SourceID = s.ID()
		chapters = append(chapters, chapter)
	})

	s.cachedChapters[manga.URL] = chapters
	return chapters, nil
}

func (s *LuaSource) PagesOf(chapter *Chapter) ([]*Page, error) {
	if cached, ok := s.cachedPages[chapter.URL]; ok {
		return cached, nil
	}

	_, err := s.call(ChapterPagesFn, lua.LTTable, lua.LString(chapter.URL))

	if err != nil {
		return nil, err
	}

	table := s.state.CheckTable(-1)
	pages := make([]*Page, 0)

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

		page.SourceID = s.ID()
		pages = append(pages, page)
	})

	s.cachedPages[chapter.URL] = pages
	return pages, nil
}

func (s *LuaSource) ID() string {
	return IDfromName(s.name)
}
