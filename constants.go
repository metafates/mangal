package main

import "time"

// Version is the current version of the program.
const Version = "2.0.1"

// Mangal is a name of application
// I keep it in a constant to avoid typos
const Mangal = "Mangal"

// CachePrefix is prefix of cache files
const CachePrefix = Mangal + "Cache"

// TempPrefix is prefix of temp files
const TempPrefix = Mangal + "Temp"

// Parallelism is number of parallel workers for scraping
const Parallelism = 100

// TestQuery is a default query for testing
const TestQuery = "Death Note"

// Forever is a constant for inifite time duration.
// It approximates to 292 years
const Forever = time.Duration(1<<63 - 1)

var AvailableFormats = []FormatType{PDF, CBZ, Plain, Zip, Epub}
var FormatsInfo = map[FormatType]string{
	PDF:   "Chapters as PDF with images",
	CBZ:   "Comic book archive format. Basically zip but with .cbz extension",
	Plain: "Just folders with raw .jpg images as chapters",
	Zip:   "Chapters compressed in zip archives",
	Epub:  "eBook format. Packs multiple chapters into single file",
}

// AsciiArt of the app
// I think it looks cool :)
const AsciiArt = "                                _\n" +
	"  /\\/\\   __ _ _ __   __ _  __ _| |\n" +
	" /    \\ / _` | '_ \\ / _` |/ _` | |\n" +
	"/ /\\/\\ \\ (_| | | | | (_| | (_| | |\n" +
	"\\/    \\/\\__,_|_| |_|\\__, |\\__,_|_|\n" +
	"                    |___/         "

// DefaultConfigStr is default config in TOML format
const DefaultConfigStr = `# Which sources to use. You can use several sources, it won't affect perfomance
use = ['manganelo']

# Type "mangal formats" to show more information about formats
format = "pdf"

# If false, then OS default reader will be used
use_custom_reader = false
custom_reader = "zathura"

# Custom download path, can be either relative (to the current directory) or absolute
download_path = '.'

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
chapter_name_template = "[%d] %s"

# Fullscreen mode 
fullscreen = true

# Input prompt symbol
prompt = ">"

# Input placeholder
placeholder = "What shall we look for?"

# Selected chapter mark
mark = "▼"

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
`
