package custom

import (
	"github.com/metafates/mangal/source"
	lua "github.com/yuin/gopher-lua"
)

func (s *LuaSource) PagesOf(chapter *source.Chapter) ([]*source.Page, error) {
	if cached, ok := s.cachedPages[chapter.URL]; ok {
		return cached, nil
	}

	_, err := s.call(ChapterPagesFn, lua.LTTable, lua.LString(chapter.URL))

	if err != nil {
		return nil, err
	}

	table := s.state.CheckTable(-1)
	pages := make([]*source.Page, 0)

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
