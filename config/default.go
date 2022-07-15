package config

var DefaultConfigBytes = []byte(`
# Which sources to use. You can use several sources, it won't affect perfomance
use = ['manganelo']

# If false, then OS default reader will be used
use_custom_reader = false
custom_reader = "zathura"




[formats]
# Type "mangal formats" to show more information about formats
default = "pdf"

# Add ComicInfo.xml to CBZ files
comicinfo = true




[downloader]
# Custom download path, can be either relative (to the current directory) or absolute
path = '.'

# How chapters should be named when downloaded
# Use %d to specify chapter number and %s to specify chapter title
# If you want to pad chapter number with zeros for natural sorting (e.g. 0001, 0123) use %0d instead of %d
chapter_name_template = "[%0d] %s"

# Add images to cache
# If set to true mangal could crash when trying to redownload something quickly
# Usually happens on slow machines
cache_images = false




[anilist]
# Enable Anilist integration (BETA)
# See https://github.com/metafates/mangal/wiki/Anilist-Integration for more information
enabled = false

# Anilist client ID
id = ""

# Anilist client secret
secret = ""

# Will mark downloaded chapters as read on Anilist
mark_downloaded = false




[ui]
# How to display chapters in TUI mode
# Use %d to specify chapter number and %s to specify chapter title
chapter_name_template = "%d %s"

# Fullscreen mode 
fullscreen = true

# Input prompt symbol
prompt = ">"

# Input placeholder
placeholder = "What shall we look for?"

# Selected chapter mark
mark = "*"

# Search window title
title = "Mangal"




[sources]
[sources.manganelo]
# Base url
base = 'https://m.manganelo.com'

# Chapters Base url
chapters_base = 'https://chap.manganelo.com/'

# Search endpoint. Put %s where the query should be
search = 'https://m.manganelo.com/search/story/%s'

# Selector of entry anchor (<a></a>) on search page
manga_anchor = '.search-story-item a.item-title'

# Selector of entry title on search page
manga_title = '.search-story-item a.item-title'

# Manga chapters anchors selector
chapter_anchor = 'li.a-h a.chapter-name'

# Manga chapters titles selector
chapter_title = 'li.a-h a.chapter-name'

# Reader page images selector
reader_page = '.container-chapter-reader img'

# Random delay between requests
random_delay_ms = 500 # ms

# Are chapters listed in reversed order on that source?
# reversed order -> from newest chapter to oldest
reversed_chapters_order = true

# With what character should the whitespace in query be replaced?
whitespace_escape = "_"
`)
