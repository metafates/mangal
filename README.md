<h1 align="center">Mangal</h1>
<p align="center">
    <img width="200" src="assets/logo.png" alt="logo">
</p>

<h3 align="center">Manga Browser & Downloader</h3>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/metafates/mangal">
    <img src="https://goreportcard.com/badge/github.com/metafates/mangal">
  </a>

  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg">
  </a>
</p>

https://user-images.githubusercontent.com/62389790/174501320-119474c3-c745-4f95-8e7d-fbf5bd40920b.mov

## Table of Contents

- [About](#about)
- [Fixing Errors](#fixing-errors)
- [Examples](#examples)
- [Config](#config)
- [Commands](#commands)
- [Install](#install)
- [Build](#build)
- [Anilist Integration](#anilist-integration)
- [Limitations](#limitations)

## About

✨ __Mangal__ is a feature rich, configurable manga browser & downloader
written in Go with support for different formats

⚙️ User defined scrapers support
 (see [config](#config) & [limitations](#limitations))

🦎 TUI & Inline modes. Use it as a standalone app or integrate with scripts

🚀 It's fast. Mangal uses multithreading to speed up the process

🍥 Integration with Anilist! __BETA__

⏳ History mode

🍿 This app is inspired by __awesome__ [ani-cli](https://github.com/pystardust/ani-cli). Check it out!


Currently, Mangal supports these formats
- PDF
- ePub
- CBZ
- ZIP
- Plain

> Type `mangal formats` for more info

## Fixing Errors

If something is not working run `mangal doctor` and follow instructions

If you still have problems,
please [open an issue](https://github.com/metafates/mangal/issues/new?assignees=metafates&labels=bug&template=bug_report.yaml) 🙏

## Examples

### TUI usage example

https://user-images.githubusercontent.com/62389790/177008719-355abc74-41c2-4e03-a7cd-ea5c2705537d.mov

https://user-images.githubusercontent.com/62389790/177008929-8c6d1c9e-a892-4479-a003-b83cfe658309.mov

### Inline mode usage example

> For more information about inline mode type `mangal inline --help`

```bash
# Search manga. Returns a list of found manga
mangal inline --query "death note"

# Search manga. Returns a JSON list of found manga
mangal inline --query "death note" --json

# Get chapters of the first manga in the list
mangal inline --query "death note" --manga 1

# Download first chapter of the first manga in the list
mangal inline --query "death note" --manga 1 --chapter 1
```

## Config

> TLDR: Use `mangal config where` to show where config should be located
> and `mangal config init` to create default config


<details>
<summary>
Config is located at the OS default config directory.
</summary>

- __Unix__ - `$XDG_CONFIG_HOME/mangal/config.toml` if `$XDG_CONFIG_HOME` exists, else `$HOME/.config/mangal/config.toml`
- __Darwin__ (macOS) - `$HOME/Library/Application\ Support/mangal/config.toml`
- __Windows__ - `%AppData%\mangal\config.toml`
</details>


You can load config from custom path by using `--config` flag or
by setting `MANGAL_CONFIG_PATH` environment variable.

> Type `mangal env` to show all supported environment variables


By default, Mangal uses [manganelo](https://m.manganelo.com/www) as a source

<details>
<summary>Click here to show config example</summary>

```toml
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
# You can use environment variable $HOME to refer to user's home directory
# If environment variable "MANGAL_DOWNLOAD_PATH" is set, then it will be used instead
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
chapter_name_template = "[%d] %s"

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
```
</details>

## Commands

```
Usage:
  mangal [flags]
  mangal [command]

Available Commands:
  cleanup     Remove cached and temp files
  completion  Generate the autocompletion script for the specified shell
  config      Config actions
  doctor      Run this in case of any errors
  env         Show environment variables
  formats     Information about available formats
  help        Help about any command
  inline      Search & Download manga in inline mode
  latest      Check if latest version of Mangal is used
  version     Show version

Flags:
  -c, --config string   use config from path
  -f, --format string   use custom format
  -h, --help            help for mangal
  -i, --incognito       do not save history
  -r, --resume          resume reading

Use "mangal [command] --help" for more information about a command.
```

## Install

- [Go (Cross Platform)](#go)
- [MacOS](#macos)
- [Windows](#windows)
- [Debian](#debian)
- [Docker](#docker)


### Go

You will need [Go installed](https://go.dev/doc/install)

```bash
go install -ldflags="-s -w" github.com/metafates/mangal@latest
```

> `-ldflags="-s -w"` - just makes the binary smaller

<details>
<summary>Update / Uninstall</summary>

#### Update

```bash
go install -ldflags="-s -w" github.com/metafates/mangal@latest
```

#### Uninstall

To uninstall just delete the binary file

- Bash / zsh - `rm $(which mangal)`
- Fish - `rm (which mangal)`
- Powershell - `rm $(where.exe mangal)`

</details>

### MacOS

Install using [Homebrew](https://brew.sh/)

```bash
brew tap metafates/mangal
brew install mangal
```

<details>
<summary>Update & Uninstall</summary>

#### Update

```bash
brew upgrade mangal
```

#### Uninstall

```bash
brew uninstall mangal
```

</details>

### Windows

Install using [Scoop](https://scoop.sh/)

```powershell
scoop install https://raw.githubusercontent.com/metafates/scoop-mangal/main/mangal.json
```

<details>
<summary>Update & Uninstall</summary>

#### Update

```powershell
scoop update mangal
```

#### Uninstall

```powershell
scoop uninstall mangal
```
</details>

### Debian

To install download the latest .deb file from [GitHub Release](https://github.com/metafates/mangal/releases) page

Then run 

```bash
sudo dpkg -i [FILE YOU DOWNLOADED].deb
```

<details>
<summary>Update & Uninstall</summary>

#### Update

To update you will need to uninstall
and install from the new .deb file in
the [GitHub Release](https://github.com/metafates/mangal/releases) page

#### Uninstall

```bash
sudo dpkg -r mangal
```

</details>

### Docker

> Thanks to @ArabCoders

Docker image is available at [Docker Hub](https://hub.docker.com/repository/docker/metafates/mangal)

You can run it by using

```bash
docker pull metafates/mangal
docker run --rm -ti -v (PWD)/mangal/downloads:/downloads -v (PWD)/mangal/config:/config metafates/mangal
```

This will create `mangal` directory in the current directory and will download manga to `mangal/downloads`

## Build

```bash
git clone https://github.com/metafates/mangal.git
cd mangal
go build -ldflags="-s -w"
```

> You can also cross build for windows, linux & macos
> by running `cross-compile.py` (you will need Python 3)
> 
> Built binaries and generated packages
> will be stored in the `bin` folder

## Anilist Integration

See [Wiki Page](https://github.com/metafates/mangal/wiki/Anilist-Integration)
for more information

## Limitations

Even though many manga sites will work,
there exist some (serious) limitations to which sites could be added

- Navigation layout should follow this model
    - Each manga have a separate page
    - Manga page should have some form of chapters list (not lazy loaded)
    - Each chapter should have a separate reader page with all the images


Some sites that work well

- https://m.manganelo.com
- https://manganato.com

See [examples of scrapers](https://github.com/metafates/mangal/discussions/7)

---

Logo taken from [here](https://www.flaticon.com/free-icon/parchment_1391306)
