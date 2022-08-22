package constant

const (
	Mangal    = "mangal"
	Version   = "3.4.1"
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"
)

const (
	TempPrefix  = Mangal + "Temp_"
	CachePrefix = Mangal + "Cache_"
)

const AssciiArtLogo = `
                                _   _____ 
  /\/\   __ _ _ __   __ _  __ _| | |___ / 
 /    \ / _' | '_ \ / _' |/ _' | |   |_ \ 
/ /\/\ \ (_| | | | | (_| | (_| | |  ___) |
\/    \/\__,_|_| |_|\__, |\__,_|_| |____/
                    |___/
`

const (
	SearchMangaFn   = "SearchManga"
	MangaChaptersFn = "MangaChapters"
	ChapterPagesFn  = "ChapterPages"
)

const SourceTemplate = `-- {{ .Name }}
-- {{ .URL }} 

-- Searches for manga with given query
-- Must return a table of tables with the following fields:
-- name: name of the manga
-- url: url of the manga
function {{ .SearchMangaFn }}(query)
end


-- Gets the list of all manga chapters
-- Returns a table of tables with the following fields:
-- name: name of the chapter
-- url: url of the chapter
function {{ .MangaChaptersFn }}(manga_url)
end


-- Gets the list of all pages of a chapter
-- Returns a table of tables with the following fields:
-- url: url of the page
-- index: index of the page
function {{ .ChapterPagesFn }}(chapter_url)
end
`
