package custom

import (
	"github.com/metafates/mangal/source"
	lua "github.com/yuin/gopher-lua"
	"strconv"
)

func (s *luaSource) ChaptersOf(manga *source.Manga) ([]*source.Chapter, error) {
	if cached, ok := s.cachedChapters[manga.URL]; ok {
		return cached, nil
	}

	_, err := s.call(mangaChaptersFn, lua.LTTable, lua.LString(manga.URL))

	if err != nil {
		return nil, err
	}

	table := s.state.CheckTable(-1)
	chapters := make([]*source.Chapter, 0)

	table.ForEach(func(k lua.LValue, v lua.LValue) {
		if k.Type() != lua.LTNumber {
			s.state.RaiseError(mangaChaptersFn + " was expected to return a table with numbers as keys, got " + k.Type().String() + " as a key")
		}

		if v.Type() != lua.LTTable {
			s.state.RaiseError(mangaChaptersFn + " was expected to return a table with tables as values, got " + v.Type().String() + " as a value")
		}

		index, err := strconv.ParseUint(k.String(), 10, 16)
		if err != nil {
			s.state.RaiseError(mangaChaptersFn + " was expected to return a table with unsigned integers as keys. " + err.Error())
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
