package constant

const (
	Mangal    = "mangal"
	Version   = "3.9.1"
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

const SourceTemplate = `{{ $divider := repeat "-" (plus (max (len .URL) (len .Name) (len .Author) 3) 12) }}{{ $divider }}
-- @name    {{ .Name }} 
-- @url     {{ .URL }}
-- @author  {{ .Author }} 
-- @license MIT
{{ $divider }}




----- IMPORTS -----
--- END IMPORTS ---




----- VARIABLES -----
--- END VARIABLES ---



----- MAIN -----

--- Searches for manga with given query.
--[[
Manga fields:
	name - string, required
 	url - string, required
	author - string, optional
	genres - string (multiple genres are divided by comma ','), optional
	summary - string, optional
--]]
-- @param query Query to search for
-- @return Table of mangas
function {{ .SearchMangaFn }}(query)
	return {}
end


--- Gets the list of all manga chapters.
--[[
Chapter fields:
	name - string, required
	url - string, required
	volume - string, optional
	manga_summary - string, optional (in case you can't get it from search page)
	manga_author - string, optional 
	manga_genres - string (multiple genres are divided by comma ','), optional
--]]
-- @param mangaURL URL of the manga
-- @return Table of chapters
function {{ .MangaChaptersFn }}(mangaURL)
	return {}
end


--- Gets the list of all pages of a chapter.
--[[
Page fields:
	url - string, required
	index - uint, required
--]]
-- @param chapterURL URL of the chapter
-- @return Table of pages
function {{ .ChapterPagesFn }}(chapterURL)
	return {}
end

--- END MAIN ---




----- HELPERS -----
--- END HELPERS ---

-- ex: ts=4 sw=4 et filetype=lua`
