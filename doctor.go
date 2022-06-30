package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func RunDoctor(vebose bool) {
	var (
		ok = func() {
			fmt.Print(successStyle.Render("OK"))
			fmt.Println()
		}

		fail = func() {
			fmt.Print(failStyle.Render("Fail"))
			fmt.Println()
		}
	)

	fmt.Print("Checking if latest version is used... ")
	latestVersion, err := FetchLatestVersion()
	if err != nil {
		fail()
		fmt.Printf("Can't find latest version\nRun %s to get more information\n", accentStyle.Render("mangal latest"))
		os.Exit(1)
	} else if latestVersion > Version {
		fail()
		fmt.Printf("New version of %s is available: %s\n", Mangal, accentStyle.Render(latestVersion))
		fmt.Printf("Run %s to get more information\n", accentStyle.Render("mangal latest"))
		os.Exit(1)
	} else {
		ok()
	}

	fmt.Print("Checking config... ")
	UserConfig = GetConfig("")

	err = ValidateConfig(UserConfig)
	if err != nil {
		fail()
		fmt.Printf("Config error: %s\n", err)
		os.Exit(1)
	}

	ok()

	var (
		sourceNotAvailable = func(source *Source) {
			fail()
			fmt.Printf("Source %s is not available\n", source.Name)
			fmt.Printf("Try to reinitialize your config with %s\n", accentStyle.Render("mangal config init --force"))
			fmt.Println("Note, that this will overwrite your current config")
			os.Exit(1)
		}

		mangaNotFound = func(source *Source, manga string) {
			query := strings.ReplaceAll(manga, " ", source.WhitespaceEscape)
			address := fmt.Sprintf(source.SearchTemplate, url.QueryEscape(strings.TrimSpace(strings.ToLower(query))))

			fail()
			fmt.Printf("Manga %s is not found\n", accentStyle.Render(manga))
			fmt.Printf("Was trying to search with %s\n\n", accentStyle.Render(address))
			fmt.Printf(
				"That probably means that %s or %s tags are invalid or website has some protection that prevents page from rendering\n",
				accentStyle.Render("manga_anchor"),
				accentStyle.Render("manga_title"),
			)

			os.Exit(1)
		}

		chaptersNotFound = func(source *Source, manga string) {
			fail()
			fmt.Printf("Chapters for %s are not found\n\n", accentStyle.Render(manga))
			fmt.Printf(
				"That probably means that %s or %s tags are invalid or website has some protection that prevents page from rendering\n",
				accentStyle.Render("chapter_anchor"),
				accentStyle.Render("chapter_title"),
			)

			os.Exit(1)
		}

		pagesNotFound = func(source *Source, manga string, chapter string) {
			fail()
			fmt.Printf("Pages for %s %s are not found\n\n", accentStyle.Render(manga), accentStyle.Render(chapter))
			fmt.Printf(
				"That probably means that %s or %s tags are invalid or website has some protection that prevents page from rendering\n",
				accentStyle.Render("page_anchor"),
				accentStyle.Render("page_title"),
			)

			os.Exit(1)
		}

		imageWasNotDownloaded = func(source *Source, manga string, chapter string, page string) {
			fail()
			fmt.Printf("Image for %s %s %s is not downloaded\n\n", accentStyle.Render(manga), accentStyle.Render(chapter), accentStyle.Render(page))
			fmt.Printf(
				"That probably means that %s tag is invalid or website has some protection that prevents page from rendering\n",
				accentStyle.Render("reader_page"),
			)

			os.Exit(1)
		}

		errorOccured = func(source *Source, action string, err error) {
			fail()
			fmt.Printf("Error occured while %s\n", action)
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	)

	scanner := bufio.NewScanner(os.Stdin)

	// check if scraper sources are online
	for _, scraper := range UserConfig.Scrapers {
		source := scraper.Source

		// read line from stdin
		fmt.Printf("Please, enter a manga title to test %s: ", source.Name)
		scanner.Scan()

		if scanner.Err() != nil {
			fail()
			fmt.Printf("Error while reading from stdin: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Checking source %s... ", source.Name)
		resp, err := http.Get(source.Base)
		if err != nil {
			sourceNotAvailable(source)
		}

		_ = resp.Body.Close()

		// check if response is 200
		if resp.StatusCode != 200 {
			sourceNotAvailable(source)
		}

		// try to get any manga page using scraper
		manga, err := scraper.SearchManga(scanner.Text())
		if err != nil {
			errorOccured(source, "searching for manga", err)
		}

		// check if manga is not empty
		if len(manga) == 0 {
			mangaNotFound(source, scanner.Text())
		}

		// get chapters for first manga
		chapters, err := scraper.GetChapters(manga[0])
		if err != nil {
			sourceNotAvailable(source)
		}

		// check if chapters is not empty
		if len(chapters) == 0 {
			chaptersNotFound(source, manga[0].Info)
		}

		// get pages for first chapter
		pages, err := scraper.GetPages(chapters[0])
		if err != nil {
			errorOccured(source, "getting pages", err)
		}

		// check if pages is not empty
		if len(pages) == 0 {
			pagesNotFound(source, manga[0].Info, chapters[0].Info)
		}

		// try to download first page
		image, err := scraper.GetFile(pages[0])
		if err != nil {
			errorOccured(source, "downloading page", err)
		}

		// check if images is not empty
		if image.Len() == 0 {
			imageWasNotDownloaded(source, manga[0].Info, chapters[0].Info, pages[0].Info)
		}

		ok()
	}
}
