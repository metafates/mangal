local goquery = require("goquery")
local http = require("http")
local client = http.client()

function SearchManga(query)
  local request = http.request("GET", "https://ww5.manganelo.tv/search/" .. query)
  local result, err = client:do_request(request)

  if err then
    error(err)
  end

  if not(result.code == 200) then
    error("code")
  end

  local doc, err = goquery.doc(result.body)

  if err ~= nil then
    error(err)
  end

  local mangas = {}

  doc:find(".item-title"):each(function (i, s)
    local manga = { name = s:text(), url = "https://ww5.manganelo.tv" .. s:attr("href") }
    mangas[i+1] = manga
  end)

  return mangas
end

function MangaChapters(manga_url)
  local request = http.request("GET", manga_url)
  local result, err = client:do_request(request)

  if err then
    error(err)
  end

  local doc, err = goquery.doc(result.body)

  if err then
    error(err)
  end

  local chapters = {}

  doc:find(".chapter-name"):each(function (i, s)
    local chapter = { name = s:text(), url = "https://ww5.manganelo.tv" .. s:attr("href") }
    chapters[i+1] = chapter
  end)

  Reverse(chapters)

  return chapters
end


function ChapterPages(chapter_url)
  local request = http.request("GET", chapter_url)
  local result, err = client:do_request(request)
  if err then
    error(err)
  end

  local doc, err = goquery.doc(result.body)

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
