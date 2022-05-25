<h1 align="center">Mangal 📖</h1>
<h3 align="center">A Manga Downloader</h3>
<p align="center">
    <img src="assets/1.jpg" alt="heh">
</p>

- [About](#about)
- [Screenshots](#screenshots)
- [Examples](#examples)
- [Install / Build](#installation--build)
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

<h3 align="center">Usage example</h4>

[![asciicast](https://asciinema.org/a/497193.svg)](https://asciinema.org/a/497193)<h3 align="center">Config example</h3>

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

<h3 align="center">Commands example</h4>

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
  version     Show version

Flags:
  -c, --config string   use config from path
  -h, --help            help for mangal

Use "mangal [command] --help" for more information about a command.
```

## Installation / Build

Currently, Mangal can be installed only by building it from source.
So you will need [go installed](https://go.dev/doc/install) to proceed further

1. `git clone https://github.com/metafates/Mangal.git`
2. `cd Mangal`
3. `make install` - `make` is used to set version string. If you can't use make (or don't want to?) feel free to just
   run `go install`.

That's it!
If you're not sure where binary is installed run `go list -f '{{.Target}}'` in the project directory

To uninstall run `make uninstall`

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

## TODO

- __Add more tests__ ⚠️
- Better error handling
- Add Mangal to package managers (homebrew, scoop, apt, ...)