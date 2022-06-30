package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func RunDoctor() {
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
		sourceNotAvailable = func(source *Source, status int) {
			fail()
			fmt.Printf("Source %s is not available\n", source.Name)
			fmt.Printf("Status code: %d\n", status)
			os.Exit(1)
		}

		mangaNotFound = func(source *Source, manga string) {
			query := strings.ReplaceAll(manga, " ", source.WhitespaceEscape)
			address := fmt.Sprintf(source.SearchTemplate, url.QueryEscape(strings.TrimSpace(strings.ToLower(query))))

			fail()

			fmt.Println(`
Manga ` + accentStyle.Render(manga) + ` was not found
Was trying to search with address: ` + accentStyle.Render(address) + `
That probably means that ` + accentStyle.Render("manga_anchor") + ` or ` + accentStyle.Render("manga_title") + ` tags are invalid, website is down or it has some protection that prevents page from rendering. 
`)

			os.Exit(1)
		}

		chaptersNotFound = func(source *Source, manga string) {
			fail()

			fmt.Println(`
Chapters for manga ` + accentStyle.Render(manga) + ` were not found
Was trying to search with address: ` + accentStyle.Render(manga) + `
That probably means that ` + accentStyle.Render("chapter_anchor") + ` or ` + accentStyle.Render("chapter_title") + ` tags are invalid or website has some protection that prevents page from rendering. 
`)

			os.Exit(1)
		}

		pagesNotFound = func(source *Source, manga string, chapter string) {
			fail()

			fmt.Println(`
Pages for chapter ` + accentStyle.Render(chapter) + ` of manga ` + accentStyle.Render(manga) + ` were not found
Was trying to search with address: ` + accentStyle.Render(manga) + `
That probably means that ` + accentStyle.Render("reader_page") + ` tag is invalid or website has some protection that prevents page from rendering. 
`)

			os.Exit(1)
		}

		imageWasNotDownloaded = func(source *Source, manga string, chapter string, page string) {
			fail()

			fmt.Println(`
Image for page ` + accentStyle.Render(page) + ` of chapter ` + accentStyle.Render(chapter) + ` of manga ` + accentStyle.Render(manga) + ` was not downloaded
Was trying to download with address: ` + accentStyle.Render(page) + `
That probably means that website has some protection that prevents image from downloading
`)
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
			errorOccured(source, "fetching base address", err)
		}

		_ = resp.Body.Close()

		// check if response is 200
		if resp.StatusCode != 200 {
			sourceNotAvailable(source, resp.StatusCode)
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
			errorOccured(source, "getting chapters", err)
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
			imageWasNotDownloaded(source, manga[0].Info, chapters[0].Info, pages[0].Address)
		}

		ok()
	}
}
