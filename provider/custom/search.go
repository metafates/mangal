package custom

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/source"
	lua "github.com/yuin/gopher-lua"
	"strconv"
)

func (s *luaSource) Search(query string) ([]*source.Manga, error) {
	if mangas := s.cache.mangas.Get(query); mangas.IsPresent() {
		m := mangas.MustGet()
		for _, manga := range m {
			manga.Source = s
		}

		return m, nil
	}

	_, err := s.call(constant.SearchMangaFn, lua.LTTable, lua.LString(query))

	if err != nil {
		return nil, err
	}

	table := s.state.CheckTable(-1)
	mangas := make([]*source.Manga, 0)

	table.ForEach(func(k lua.LValue, v lua.LValue) {
		if k.Type() != lua.LTNumber {
			s.state.RaiseError(constant.SearchMangaFn + " was expected to return a table with numbers as keys, got " + k.Type().String() + " as a key")
		}

		if v.Type() != lua.LTTable {
			s.state.RaiseError(constant.SearchMangaFn + " was expected to return a table with tables as values, got " + v.Type().String() + " as a value")
		}

		index, err := strconv.ParseUint(k.String(), 10, 16)
		if err != nil {
			s.state.RaiseError(constant.SearchMangaFn + " was expected to return a table with unsigned integers as keys. " + err.Error())
		}

		manga, err := mangaFromTable(v.(*lua.LTable), uint16(index))

		if err != nil {
			s.state.RaiseError(err.Error())
		}

		manga.Source = s
		mangas = append(mangas, manga)
	})

	_ = s.cache.mangas.Set(query, mangas)
	return mangas, nil
}
