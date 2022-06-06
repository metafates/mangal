<h1 align="center">Mangal 📖</h1>
<p align="center">
    <img width="200" src="assets/manga.png" alt="heh">
</p>

<h3 align="center">Manga Downloader</h3>

- [About](#about)
- [Screenshots](#screenshots)
- [Examples](#examples)
- [Install](#install)
- [Build](#build)
- [Limitations](#limitations)

## About

✨ __Mangal__ is a fancy TUI app written in go that scrapes, downloads and packs manga into pdfs

⚙️ The most important feature of Mangal is that it supports user defined scrapers
that can be added with just a few lines of config file (see [limitations](#limitations))

🧋 Built with [Bubble Tea framework](https://github.com/charmbracelet/bubbletea)

🍿 This app is inspired by __awesome__ [ani-cli](https://github.com/pystardust/ani-cli). Check it out!

## Screenshots

<img alt="search input" src="assets/sc1.png">
<img alt="list with found manga" src="assets/sc2.png">
<img alt="chapters of choosen manga" src="assets/sc3.png">
<img alt="prompt before downloading" src="assets/sc4.png">
<img alt="downloading progress bar" src="assets/sc5.png">

## Examples

<h3 align="center">Usage example</h3>

https://user-images.githubusercontent.com/62389790/172178993-bb40392a-54ba-446d-b0ed-1b962ede7ed2.mp4


<h3 align="center">Config example</h3>

Config is located at the OS default config directory.

- __Unix__ - `$XDG_CONFIG_HOME/mangal/config.toml` if `$XDG_CONFIG_HOME` exists, else `$HOME/.config/mangal/config.toml`
- __Darwin__ (macOS) - `$HOME/Library/Application\ Support/mangal/config.toml`
- __Windows__ - `%AppData%\mangal\config.toml`
- __Plan 9__ - `$home/lib/mangal/config.toml`

You can load config from custom path by using `--config` flag

`mangal --config /user/configs/config.toml`

By default, Mangal uses [manganelo](https://ww5.manganelo.tv) as a source

```toml
# Which sources to use. You can use several sources in descendant order priority
use = ['manganelo']

# Default download path
path = '.'

# Fullscreen mode
fullscreen = true

[sources]
    [sources.manganelo]
    # Base url
    base = 'https://ww5.manganelo.tv'

    # Search endpoint. Put %s where the query should be
    search = 'https://ww5.manganelo.tv/search/%s'

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
    
    # Random delay between requests (in miliseconds)
    random_delay_ms = 500
    
    # Are chapters listed in reversed order on that source?
    # reversed = from latest to oldest
    reversed_chapters_order = true

```

<h3 align="center">Commands example</h3>

```
$ mangal help

A fast and flexible manga downloader

Usage:
  mangal [flags]
  mangal [command]

Available Commands:
  cleanup     Remove cached and temp files
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Create default config at default path
  update      Update Mangal
  version     Show version
  where       Show path where config is located

Flags:
  -c, --config string   use config from path
  -h, --help            help for mangal

Use "mangal [command] --help" for more information about a command.
```

## Install

### Homebrew

```bash
brew tap metafates/tap
brew install metafates/tap/mangal
```

### Go
```bash
go install github.com/metafates/mangal@latest
```

## Build

```bash
git clone https://github.com/metafates/mangal.git
cd mangal
go build
```

## Limitations

Even though most manga sites will work, there exists some limitation to which sites could be added

- Navigation layout should follow this model
    - Each manga have a separate page
    - Manga page should have a some form of chapters list (not lazy loaded)
    - Each chapter should have a separate page with pages (images)

<br>

Some sites that work well

- https://manganato.com
- https://ww3.mangakakalot.tv
- https://ww5.manganelo.tv

---

Manga image taken from [here](https://www.flaticon.com/free-icons/manga)
