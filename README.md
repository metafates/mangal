<h1 align="center">Mangal 3 📜</h1>

<p align="center">
    <img alt="Linux" src="https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black">
    <img alt="macOS" src="https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=macos&logoColor=F0F0F0">
    <img alt="Windows" src="https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white">
</p>

<h3 align="center">
    The most advanced CLI manga downloader in the entire universe!
</h3>

https://user-images.githubusercontent.com/62389790/183284495-86140f8b-d543-4bc4-a413-37cb07c1552e.mov


## Table of contents
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Custom scrapers](#custom-scrapers)
- [Anilist](#anilist)

## Features

- __LUAAAA SCRAPPEERRRSS!!!__ You can add any source you want by creating your own _(or using someone's else)_ scraper with __Lua 5.1__. See [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers)
- __Download & Read Manga__ - I mean, it would be strange if you couldn't, right?
- __4 Different export formats__ - PDF, CBZ, ZIP and plain images
- __Fast__ - yes. 
- __Monolith__ - ZERO runtime dependencies. Even Lua is built in.
- __Fancy__ - (ﾉ>ω<)ﾉ :｡･:*:･ﾟ’★,｡･:*:･ﾟ’☆
- __Cross-Platform__ - Linux, macOS, Windows. Sorry, FreeBSD users...
- __Anilist integration__ - Track your manga progress on Anilist when reading with Mangal.

## Installation

### Go (Any OS)

Visit this link to install [Go](https://go.dev/doc/install)

    go install -ldflags="-s -w" github.com/metafates/mangal@latest

> **Info** `-ldflags="-s -w"` makes the binary smaller
> 
> Use this method if others are not working for some reason.
> And please open an issue if so


### Linux

[AUR package](https://aur.archlinux.org/packages/mangal-bin) (by @balajsra)

Download the latest version from [GitHub release page](https://github.com/metafates/mangal/releases/latest)

### macOS

Install using [Homebrew](https://brew.sh/)

    brew tap metafates/mangal
    brew install mangal

### Windows

Install using [Scoop](https://scoop.sh/)

    scoop bucket add metafates https://github.com/metafates/scoop-metafates
    scoop install mangal

### Docker

Install using... well, you know.

    docker pull metafates/mangal

To run

    docker run --rm -ti -e "TERM=xterm-256color" -v (PWD)/mangal/downloads:/downloads -v (PWD)/mangal/config:/config metafates/mangal

## Usage

### TUI

Just run `mangal` and you're ready to go.

### Mini

There's also a `mini` mode that tries to mimic [ani-cli](https://github.com/pystardust/ani-cli)

To run: `mangal mini`

<img width="254" alt="Screenshot 2022-08-14 at 09 37 14" src="https://user-images.githubusercontent.com/62389790/184524070-88fd36f7-9875-4a41-904c-04caad110549.png">

### Other

See `mangal help` for more information

## Configuration

Mangal uses [TOML](https://toml.io) format for configuration under the `mangal.toml` filename.
Config is expected to be at the OS default config directory.
For example, on Linux it would be `~/.config/mangal/mangal.toml`.

Run `mangal where` to show expected config paths

> "But what if I want to specify my own config path?"
> 
> Okay, fine, use env variable `MANGAL_CONFIG_PATH`

Run `mangal config init` to generate a default config file

> This is not a complete config, just an example.
```toml
# mangal.toml

[downloader]
# Name template of the downloaded chapters
# Available variables:
# {index}        - index of the chapters
# {padded-index} - same as index but padded with leading zeros
# {chapter}      - name of the chapter
chapter_name_template = '[{padded-index}] {chapter}'

# Where to download manga
# Absolute or relative.
#
# You can also use home variable 
# Linux/macOS = $HOME or ~
# Windows = %USERPROFILE%
path = '.'

# Use asynchronous downloader (faster)
# Do no turn it off unless you have some issues
async = true



[formats]
# Default format to export chapters
# Available options are: pdf, zip, cbz, plain
use = 'pdf'



[history]
# Save chapters to history when downloaded
save_on_download = false
# Save chapters to history on read
save_on_read = true



[icons]
# Icons variant.
# Available options are: emoji, kaomoji, plain, squares, nerd (nerd-font)
variant = 'emoji'



[mangadex]
# Preffered language
language = 'en'
# Show nsfw manga/chapters
nsfw = false
# Show chapters that cannot be read (because they are hosted on a different site)
show_unavailable_chapters = false



[mini]
# Limit number of search page entries
search_limit = 20



[reader]
# Name of the app to use as a reader. Will use default OS app if empty
name = ''
# Will open chapter in the browser instead of downloading it
read_in_browser = false

[logs]
# write logs?
write = false
# Available options are: (from less to most verbose)
# panic, fatal, error, warn, info, debug, trace
level = "info"
```

## Custom scrapers

Mangal has a Lua5.1 VM built-in + some useful libraries, such as headless chrome, http client, html parser and so on...

Check the [defined modules](luamodules) for more information.

For scraper examples, check the [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers)

_Okay, so, how do I add a custom scraper?_

1. Create a new lua file in the `mangal where --sources` folder
2. Filename will be used as a source name
3. Your script __must__ contain __3__ essential functions
   - `SearchManga(query)` - must return a table of tables each having 2 fields `name` and `url`
   - `MangaChapters(mangalUrl)` - must return a table of tables each having 2 fields `name` and `url` _(again)_
   - `ChapterPages(chapterUrl)` - must return a table of tables each having 2 fields `index` _(for ordering)_ and `url` _(to download image)_
4. __That's it!__ You can test it by running `mangal run ...` where `...` is a filename

New to Lua? [Quick start guide](https://learnxinyminutes.com/docs/lua/)

## Anilist

Mangal also supports integration with anilist.

It will mark chapters as read on Anilsit when you read them inside mangal.

For more information see [wiki](https://github.com/metafates/mangal/wiki/Anilist-Integration)

> Maybe I'll add more sites in the future, like [myanimelist](https://myanimelist.net/). Open to suggestions!
