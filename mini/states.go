package mini

import (
	"fmt"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/converter"
	"github.com/metafates/mangal/provider"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/style"
	"github.com/samber/lo"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"regexp"
	"strconv"
	"strings"
)

type state int

const (
	mangasSearchState state = iota + 1
	mangaSelectState
	sourceSelectState
	chapterSelectState
	chapterReadState
	chaptersDownloadState
	quitState
)

func (m *mini) handleSourceSelectState() error {
	defaultProviders := provider.DefaultProviders()
	customProviders := lo.Must(provider.CustomProviders())

	var providers = make([]*provider.Provider, 0)

	for _, p := range defaultProviders {
		providers = append(providers, p)
	}

	for _, p := range customProviders {
		providers = append(providers, p)
	}

	slices.SortFunc(providers, func(a *provider.Provider, b *provider.Provider) bool {
		return strings.Compare(a.String(), b.String()) < 0
	})

	var err error

	title("Select Source")
	b, p, err := menu(providers)
	if err != nil {
		return err
	}

	if quit.eq(b) {
		m.newState(quitState)
		return nil
	}

	erase := printErasable(style.Blue("Initializing Source..."))
	m.selectedSource, err = p.CreateSource()
	erase()

	m.newState(mangasSearchState)
	return err
}

func (m *mini) handleMangaSearchState() error {
	var searchLoop func() error
	title("Search Manga")

	searchLoop = func() error {
		in, err := getInput(func(s string) bool {
			return s != ""
		})

		if err != nil {
			return err
		}

		query := in.value

		erase := printErasable(style.Blue("Searching Query..."))
		m.cachedMangas[query], err = m.selectedSource.Search(query)
		erase()

		if len(m.cachedMangas[query]) == 0 {
			fail("No search results found")
			return searchLoop()
		}

		m.query = query
		m.newState(mangaSelectState)
		return err
	}

	return searchLoop()
}

func (m *mini) handleMangaSelectState() error {
	var err error
	title("Query Results >>")
	b, p, err := menu(m.cachedMangas[m.query])
	if err != nil {
		return err
	}

	if quit.eq(b) {
		m.newState(quitState)
		return nil
	}

	m.selectedManga = p
	m.newState(chapterSelectState)
	return err
}

func (m *mini) handleChapterSelectState() error {
	var err error

	erase := printErasable(style.Magenta("Searching Chapters..."))
	m.cachedChapters[m.selectedManga.URL], err = m.selectedSource.ChaptersOf(m.selectedManga)
	erase()
	if err != nil {
		return err
	}

	chapters := m.cachedChapters[m.selectedManga.URL]

	if len(chapters) == 0 {
		fail("No chapters found")
		m.selectedManga = nil
		m.newState(mangaSelectState)
		return nil
	}

	title(fmt.Sprintf("To specify a range, use: start_number end_number (Episodes: 1-%d)", len(chapters)))
	oneChapterInput := regexp.MustCompile(`^\d+$`)
	rangeInput := regexp.MustCompile(`^\d+ \d+$`)
	in, err := getInput(func(s string) bool {
		var err error

		switch {
		case rangeInput.MatchString(s):
			var a, b int64
			{
				l := strings.Split(s, " ")
				a, err = strconv.ParseInt(l[0], 10, 16)
				if err != nil {
					return false
				}

				b, err = strconv.ParseInt(l[1], 10, 16)
				if err != nil {
					return false
				}
			}

			return a < b && 0 < a && int(a) < len(chapters) && int(b) <= len(chapters)
		case oneChapterInput.MatchString(s):
			var a int64
			a, err = strconv.ParseInt(s, 10, 16)
			if err != nil {
				return false
			}

			return 0 < a && int(a) <= len(chapters)
		default:
			return false
		}
	})

	if err != nil {
		return err
	}

	switch {
	case rangeInput.MatchString(in.value):
		nums := strings.Split(in.value, " ")
		from := lo.Must(strconv.ParseInt(nums[0], 10, 16))
		to := lo.Must(strconv.ParseInt(nums[1], 10, 16))

		for i := from - 1; i < to; i++ {
			m.selectedChapters = append(m.selectedChapters, chapters[i])
		}
	case oneChapterInput.MatchString(in.value):
		num := lo.Must(strconv.ParseInt(in.value, 10, 16))
		m.selectedChapters = append(m.selectedChapters, chapters[num-1])
	}

	if m.download {
		m.newState(chaptersDownloadState)
	} else {
		m.newState(chapterReadState)
	}

	return nil
}

func (m *mini) handleChapterReadState() error {
	var (
		err      error
		readLoop func(*source.Chapter) (bool, error)
	)

	readLoop = func(chapter *source.Chapter) (bool, error) {
		erase := printErasable(style.Blue("Loading Chapter..."))
		m.cachedPages[chapter.URL], err = m.selectedSource.PagesOf(chapter)
		erase()
		if err != nil {
			return false, err
		}

		erase = printErasable(style.Blue("Downloading Pages..."))
		err = chapter.DownloadPages()
		erase()

		if err != nil {
			return false, err
		}

		erase = printErasable(style.Blue("Converting..."))
		conv, err := converter.Get(viper.GetString(config.FormatsUse))
		if err != nil {
			return false, err
		}

		path, err := conv.SaveTemp(chapter)
		erase()

		if err != nil {
			return false, err
		}

		erase = printErasable(style.Blue("Opening..."))

		if reader := viper.GetString(config.ReaderName); reader != "" {
			err = open.RunWith(path, reader)
			if err != nil {
				return false, err
			}
		} else {
			err = open.Run(path)
			if err != nil {
				return false, err
			}
		}

		erase()

		title(fmt.Sprintf("Currently reading %s", chapter.Name))
		b, _, err := menu([]*source.Chapter{}, &next, &reread, &back)
		if err != nil {
			return false, err
		}

		switch b {
		case &next:
			return false, nil
		case &reread:
			return readLoop(chapter)
		case &back:
			m.newState(chapterSelectState)
			return false, nil
		case &quit:
			m.newState(quitState)
			return true, nil
		default:
			return false, nil
		}
	}

	for _, chapter := range m.selectedChapters {
		q, err := readLoop(chapter)
		if err != nil {
			return err
		}

		if q {
			return nil
		}
	}

	return nil
}
