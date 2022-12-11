package constant

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


---@alias manga { name: string, url: string, author: string|nil, genres: string|nil, summary: string|nil }
---@alias chapter { name: string, url: string, volume: string|nil, manga_summary: string|nil, manga_author: string|nil, manga_genres: string|nil }
---@alias page { url: string, index: number }


----- IMPORTS -----
--- END IMPORTS ---




----- VARIABLES -----
--- END VARIABLES ---



----- MAIN -----

--- Searches for manga with given query.
-- @param query string Query to search for
-- @return manga[] Table of mangas
function {{ .SearchMangaFn }}(query)
	return {}
end


--- Gets the list of all manga chapters.
-- @param mangaURL string URL of the manga
-- @return chapter[] Table of chapters
function {{ .MangaChaptersFn }}(mangaURL)
	return {}
end


--- Gets the list of all pages of a chapter.
-- @param chapterURL string URL of the chapter
-- @return page[]
function {{ .ChapterPagesFn }}(chapterURL)
	return {}
end

--- END MAIN ---




----- HELPERS -----
--- END HELPERS ---

-- ex: ts=4 sw=4 et filetype=lua
`
