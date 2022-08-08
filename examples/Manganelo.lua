local html = require("html")
local http = require("http")
local client = http.client()

local manganelo = "https://ww5.manganelo.tv"


function SearchManga(query)
  local request = http.request("GET", manganelo .. "/search/" .. query)
  local result, err = client:do_request(request)
  AssertFalse(err)

  local doc = html.parse(result.body)
  local mangas = {}

  doc:find(".item-title"):each(function (i, s)
    local manga = { name = s:text(), url = manganelo .. s:attr("href") }
    mangas[i+1] = manga
  end)

  return mangas
end


function MangaChapters(manga_url)
  local request = http.request("GET", manga_url)
  local result, err = client:do_request(request)
  AssertFalse(err)

  local doc = html.parse(result.body)

  local chapters = {}

  doc:find(".chapter-name"):each(function (i, s)
    local chapter = { name = s:text(), url = manganelo .. s:attr("href") }
    chapters[i+1] = chapter
  end)

  Reverse(chapters)

  return chapters
end


function ChapterPages(chapter_url)
  local request = http.request("GET", chapter_url)
  local result, err = client:do_request(request)
  AssertFalse(err)

  local doc = html.parse(result.body)

  local pages = {}

  doc:find(".container-chapter-reader img"):each(function (i, s)
    local page = { index = i, url = s:attr("data-src") }
    pages[i+1] = page
  end)

  return pages
end


function Reverse(t)
  local n = #t
  local i = 1
  while i < n do
    t[i],t[n] = t[n],t[i]
    i = i + 1
    n = n - 1
  end
end


function AssertFalse(e)
  if e then error(e) end
end

