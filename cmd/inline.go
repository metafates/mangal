package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metafates/mangal/cleaner"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/downloader"
	"github.com/metafates/mangal/scraper"
	"github.com/metafates/mangal/util"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"os"
	"sync"
)

// inlineOptions provides all options for inline mode
type inlineOptions struct {
	mangaIdx   int
	chapterIdx int
	asJson     bool
	format     common.FormatType
	showUrls   bool
	asTemp     bool
	doRead     bool
	doOpen     bool
}

// inlineMode provides all functionality of TUI but in inline mode
// TODO: split into subfunctions
func inlineMode(query string, options inlineOptions) (string, error) {
	if !options.asTemp {
		defer cleaner.RemoveTemp()
	}

	if options.format != "" {
		config.UserConfig.Formats.Default = options.format
	}

	// Check if config is valid
	if err := config.ValidateConfig(config.UserConfig); err != nil {
		return "", err
	}

	var (
		manga []*scraper.URL
		wg    sync.WaitGroup
	)

	// Check if query is valid
	if query == "" {
		return "", errors.New("query expected")
	}

	wg.Add(len(config.UserConfig.Scrapers))

	// Search for manga in all scrapers
	for _, s := range config.UserConfig.Scrapers {
		go func(s *scraper.Scraper) {
			defer wg.Done()

			m, err := s.SearchManga(query)

			if err == nil {
				manga = append(manga, m...)
			}
		}(s)
	}

	wg.Wait()

	// Check if manga was selected
	if options.mangaIdx >= 0 {
		if options.mangaIdx > len(manga) || options.mangaIdx <= 0 {
			return "", errors.New("index out of range")
		}

		selectedManga := manga[options.mangaIdx-1]

		// Get chapters of selected manga
		chapters, err := selectedManga.Scraper.GetChapters(selectedManga)
		if err != nil {
			return "", errors.New("error while getting chapters")
		}

		// Check if chapter was selected
		if options.chapterIdx >= 0 {

			// Get selected chapter
			selectedChapter, ok := util.Find(chapters, func(c *scraper.URL) bool {
				return c.Index == options.chapterIdx
			})

			if !ok {
				return "", errors.New("index out of range")
			}

			// if option to open chapter is set, open it
			if options.doOpen {
				if err = open.Start(selectedChapter.Address); err != nil {
					return "", errors.New("unexpected error while trying to open url")
				}

				return "", nil
			}

			// if option to read chapter is set download chapter as temp file
			if options.doRead {
				options.asTemp = true
			}

			// Download chapter
			chapterPath, err := downloader.DownloadChapter(selectedChapter, nil, options.asTemp)
			if err != nil {
				return "", errors.New("error while downloading chapter")
			}

			// if epub file was used create it
			if downloader.EpubFile != nil {
				downloader.EpubFile.SetAuthor(selectedManga.Scraper.Source.Base)
				if err := downloader.EpubFile.Write(chapterPath); err != nil {
					return "", errors.New("error while making epub file")
				}

				// reset epub file
				downloader.EpubFile = nil
			}

			// if options to read chapter is set, read it
			if options.doRead {
				// check if custom reader is set
				if config.UserConfig.Reader.UseCustomReader {
					err = open.StartWith(chapterPath, config.UserConfig.Reader.CustomReader)
				} else {
					err = open.Start(chapterPath)
				}

				if err != nil {
					return "", err
				}

				return "", nil
			}

			return chapterPath, nil
		}

		// if option to print data as json is set, print it as json
		if options.asJson {
			data, err := json.Marshal(chapters)
			if err != nil {
				return "", errors.New("could not get data as json")
			}

			return string(data), nil
		} else if options.doOpen {
			if err = open.Start(selectedManga.Address); err != nil {
				return "", errors.New("unexpected error while trying to open url")
			}
			return "", nil
		} else {
			var chaptersString string

			// print chapters list
			for _, c := range chapters {
				if options.showUrls {
					chaptersString += fmt.Sprintf("[%d] %s %s\n", c.Index, c.Info, c.Address)
				} else {
					chaptersString += fmt.Sprintf("[%d] %s\n", c.Index, c.Info)
				}
			}

			return chaptersString, nil
		}

	} else {
		// if option to print data as json is set, print it as json
		if options.asJson {
			data, err := json.Marshal(manga)
			if err != nil {
				return "", errors.New("could not get data as json")
			}

			return string(data), nil
		} else {
			var mangaString string

			// print manga list
			for i, m := range manga {
				if options.showUrls {
					mangaString += fmt.Sprintf("[%d] %s %s\n", i+1, m.Info, m.Address)
				} else {
					mangaString += fmt.Sprintf("[%d] %s\n", i+1, m.Info)
				}
			}

			return mangaString, nil
		}
	}
}

var inlineCmd = &cobra.Command{
	Use:   "inline",
	Short: "Search & Download manga in inline mode",
	Long: `Search & Download manga in inline mode
Useful for scripting`,
	Run: func(cmd *cobra.Command, args []string) {
		query, _ := cmd.Flags().GetString("query")
		mangaIdx, _ := cmd.Flags().GetInt("manga")
		chapterIdx, _ := cmd.Flags().GetInt("chapter")
		asJson, _ := cmd.Flags().GetBool("json")
		format, _ := cmd.Flags().GetString("format")
		showUrls, _ := cmd.Flags().GetBool("urls")
		asTemp, _ := cmd.Flags().GetBool("temp")
		doRead, _ := cmd.Flags().GetBool("read")
		doOpen, _ := cmd.Flags().GetBool("open")

		res, err := inlineMode(query, inlineOptions{
			mangaIdx:   mangaIdx,
			chapterIdx: chapterIdx,
			asJson:     asJson,
			format:     common.FormatType(format),
			showUrls:   showUrls,
			asTemp:     asTemp,
			doRead:     doRead,
			doOpen:     doOpen,
		})

		if err != nil {
			if asJson {
				fmt.Printf(`{error: "%s"}\n`, err)
			} else {
				fmt.Println(err)
			}

			os.Exit(1)
		}

		fmt.Println(res)
	},
}

func init() {
	inlineCmd.Flags().Int("manga", -1, "choose manga by index")
	inlineCmd.Flags().Int("chapter", -1, "choose and download chapter by index")
	inlineCmd.Flags().StringP("query", "q", "", "manga to search")
	inlineCmd.Flags().BoolP("json", "j", false, "print as json")
	inlineCmd.Flags().StringP("format", "f", "", "use custom format")
	inlineCmd.Flags().BoolP("urls", "u", false, "show urls")
	inlineCmd.Flags().BoolP("temp", "t", false, "download as temp")
	inlineCmd.Flags().BoolP("read", "r", false, "read chapter")
	inlineCmd.Flags().BoolP("open", "o", false, "open url")
	inlineCmd.Flags().SortFlags = false
	_ = inlineCmd.MarkFlagRequired("query")
	_ = inlineCmd.MarkFlagFilename("config", "toml")
	_ = inlineCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return util.Map(common.AvailableFormats, util.ToString[common.FormatType]), cobra.ShellCompDirectiveDefault
	})
	mangalCmd.AddCommand(inlineCmd)
}
