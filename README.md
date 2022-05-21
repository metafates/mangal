<h1 align="center">Mangai 📖</h1>
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

✨ __Mangai__ is a fancy TUI app written in go that scrapes, downloads and packs manga into pdfs


⚙️ The most important feature of Mangai is that it supports user defined scrapers
that can be added with just a few lines of config file (see [limitations](#limitations))

🧋 Built with [Bubble Tea framework](https://github.com/charmbracelet/bubbletea)

🍿 This app is inspired by __awesome__ [ani-cli](https://github.com/pystardust/ani-cli). Check it out!

## Screenshots

<img src="assets/sc1.png">
<img src="assets/sc2.png">
<img src="assets/sc3.png">
<img src="assets/sc4.png">

## Examples

<h3 align="center">Usage example</h4>

[![asciicast](https://asciinema.org/a/Kr4xdcfndSdvQCoWpoBNIyUFH.svg)](https://asciinema.org/a/Kr4xdcfndSdvQCoWpoBNIyUFH)
<br><br>
<h3 align="center">Config example</h3>

Config is located at the OS default config directory.

- __Unix__ - `$XDG_CONFIG_HOME/mangai/config.toml` if `$XDG_CONFIG_HOME` exists, else `$HOME/.config/mangai/config.toml`
- __Darwin__ (macOS) - `$HOME/Library/Application\ Support/mangai/config.toml`
- __Windows__ - `%AppData%\mangai\config.toml`
- __Plan 9__ - `$home/lib/mangai/config.toml`

Custom config paths not supported _(yet)_

_By default (if no config defined) Mangai uses [manganelo](https://ww5.manganelo.tv) as a source_

```toml
# Which sources to use for searching.
# Since searching is done asynchronously it should not affect perfomance
use = ['manganelo', 'mangapoisk']

# Default download path
# It could be relative or absolute 
# path = '/users/user/manga'
path = '.'

# Fullscreen mode
fullscreen = true

# This is where you define new sources for searching
[sources]
    # sources.%name of the source%
    [sources.manganelo]
    
    # Base url of the source
    base = 'https://ww5.manganelo.tv'

    # Search endpoint. Put %s where the query should be
    search = 'https://ww5.manganelo.tv/search/%s'

    # Selector of entry anchor (<a></a>) on search page
    # Make this and other selectors as specific as possible
    # See https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Selectors
    manga_anchor = '.search-story-item a.item-title'

    # Selector of entry title on search page
    manga_title = '.search-story-item a.item-title'

    # Manga chapters anchors selector
    chapter_anchor = 'li.a-h a.chapter-name'

    # Manga chapters titles selector
    chapter_title = 'li.a-h a.chapter-name'

    # Reader page images selector
    chapter_panels = '.container-chapter-reader img'

    [sources.mangapoisk]
    base = 'https://mangapoisk.ru'
    search = 'https://mangapoisk.ru/search?q=%s'
    manga_anchor = 'article.card a.px-1'
    manga_title = 'article.card .entry-title'
    chapter_anchor = 'li.chapter-item a'
    chapter_title = 'li.chapter-item span.chapter-title'
    chapter_panels = 'div.chapter-container .chapter-images img'
```

## Installation / Build

Currently, Mangai can be installed only by building it from source.
So you will need [go installed](https://go.dev/doc/install) to proceed further

1. `git clone https://github.com/metafates/Mangai.git`
2. `cd mangai`
3. `make install` - `make` is used to set version string. If you can't use make (or don't want to?) feel free to just run `go install`.

That's it!
If you're not sure where binary is installed run `go list -f '{{.Target}}'` in the project directory

To uninstall run `make uninstall`

## Limitations

Even though most manga sites will work, there exists some limitation to which sites could be added

- Navigation layout should follow this model
    - Each manga have a separate page
    - Manga page should have a some form of chapters list (not lazy loaded)
    - Each chapter should have a separate page with panels (images)
- No anti-bot protection 🤖

<br>

Some sites that work well
- https://manganato.com
- https://ww3.mangakakalot.tv
- https://ww5.manganelo.tv
- https://mangapoisk.ru
