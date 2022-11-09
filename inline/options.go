package inline

import (
	"fmt"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type (
	MangaPicker    func([]*source.Manga) *source.Manga
	ChaptersFilter func([]*source.Chapter) ([]*source.Chapter, error)
)

type Options struct {
	Out                 io.Writer
	Sources             []source.Source
	IncludeAnilistManga bool
	Download            bool
	Json                bool
	PopulatePages       bool
	Query               string
	MangaPicker         mo.Option[MangaPicker]
	ChaptersFilter      mo.Option[ChaptersFilter]
}

func ParseMangaPicker(query, description string) (MangaPicker, error) {
	const (
		first = "first"
		last  = "last"
		exact = "exact"
	)

	pattern := fmt.Sprintf(`^(%s|%s|%s|\d+)$`, first, last, exact)
	mangaPickerRegex := regexp.MustCompile(pattern)

	if !mangaPickerRegex.MatchString(description) {
		return nil, fmt.Errorf("invalid manga picker pattern: %s", description)
	}

	return func(mangas []*source.Manga) *source.Manga {
		if len(mangas) == 0 {
			return nil
		}

		switch description {
		case first:
			return mangas[0]
		case last:
			return mangas[len(mangas)-1]
		case exact:
			for _, manga := range mangas {
				if manga.Name == query {
					return manga
				}
			}

			return nil
		default:
			index := lo.Must(strconv.ParseUint(description, 10, 16))
			return mangas[util.Min(index, uint64(len(mangas)-1))]
		}
	}, nil
}

func ParseChaptersFilter(description string) (ChaptersFilter, error) {
	const (
		first = "first"
		last  = "last"
		all   = "all"
		from  = "From"
		to    = "To"
		sub   = "Sub"
	)

	pattern := fmt.Sprintf(`^(%s|%s|%s|(?P<%s>\d+)(-(?P<%s>\d+))?|@(?P<%s>.+)@)$`, first, last, all, from, to, sub)
	mangaPickerRegex := regexp.MustCompile(pattern)

	if !mangaPickerRegex.MatchString(description) {
		return nil, fmt.Errorf("invalid chapter filter pattern: %s", description)
	}

	return func(chapters []*source.Chapter) ([]*source.Chapter, error) {
		if len(chapters) == 0 {
			return chapters, nil
		}

		switch description {
		case first:
			return chapters[0:1], nil
		case last:
			return chapters[len(chapters)-1:], nil
		case all:
			return chapters, nil
		default:
			groups := util.ReGroups(mangaPickerRegex, description)

			if sub, ok := groups[sub]; ok && sub != "" {
				return lo.Filter(chapters, func(a *source.Chapter, _ int) bool {
					return strings.Contains(a.Name, sub)
				}), nil
			}

			from := lo.Must(strconv.ParseUint(groups[from], 10, 16))
			from = util.Min(from, uint64(len(chapters)))

			n := groups[to]
			if n == "" {
				return []*source.Chapter{chapters[from]}, nil
			}

			to := lo.Must(strconv.ParseUint(n, 10, 16))
			to = util.Min(to, uint64(len(chapters)))

			if from > to {
				from, to = to, from
			}

			return chapters[from : to+1], nil
		}
	}, nil
}
