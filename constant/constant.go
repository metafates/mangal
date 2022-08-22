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

const SourceTemplate = `{{ repeat "-" (max (len .URL) (len .Name)) }}{{ repeat "-" 12 }}
-- @name    {{ .Name }} 
-- @url     {{ .URL }}
-- @author  {{ .Author }} 
-- @license MIT
{{ repeat "-" (max (len .URL) (len .Name)) }}{{ repeat "-" 12 }}


--- IMPORTS ---
-- ...


--- VARIABLES ---
-- ...


--- Searches for manga with given query.
-- @param query Query to search for
-- @return Table of tables with the following fields: name, url
function {{ .SearchMangaFn }}(query)
	return {}
end


--- Gets the list of all manga chapters.
-- @param mangaURL URL of the manga
-- @return Table of tables with the following fields: name, url
function {{ .MangaChaptersFn }}(mangaURL)
	return {}
end


--- Gets the list of all pages of a chapter.
-- @param chapterURL URL of the chapter
-- @return Table of tables with the following fields: url, index
function {{ .ChapterPagesFn }}(chapterURL)
	return {}
end


-- ex: ts=4 sw=4 et filetype=lua`
