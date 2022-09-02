<h1 align="center">Mangal 3 📜</h1>

<p align="center">
    <img alt="Linux" src="https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black">
    <img alt="macOS" src="https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=macos&logoColor=F0F0F0">
    <img alt="Windows" src="https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white">
    <img alt="Termux" src="https://img.shields.io/badge/Termux-000000?style=for-the-badge&logo=GNOME%20Terminal&logoColor=white">
</p>

<h3 align="center">
    The most advanced CLI manga downloader in the entire universe!
</h3>

https://user-images.githubusercontent.com/62389790/183284495-86140f8b-d543-4bc4-a413-37cb07c1552e.mov

## Try it!

```shell
curl -sfL io.metafates.one/mr  | sh
```

> **Note** This script does not install anything, it just downloads, verifies and runs Mangal.
> Not available on Windows.

## Table of contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Custom scrapers](#custom-scrapers)
- [Anilist](#anilist)

## Features

- __Lua Scrapers!!!__ You can add any source you want by creating your own _(or using someone's else)_ scraper with __
  Lua 5.1__. See [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers)
- [Mangadex](https://mangadex.org) + [Manganelo](https://m.manganelo.com/wwww) Built-In
- __Download & Read Manga__ - I mean, it would be strange if you couldn't, right?
- __4 Different export formats__ - PDF, CBZ, ZIP and plain images
- __3 Different modes__ - TUI, Mini and Inline
- __Fast?__ - YES.
- __Monolith__ - ZERO runtime dependencies. Even Lua is built in.
- __Fancy__ - (ﾉ>ω<)ﾉ :｡･::･ﾟ’★,｡･:･ﾟ’☆
- __Cross-Platform__ - Linux, macOS, Windows, Termux (partially)
- __Anilist integration__ - Track your manga progress on Anilist when reading with Mangal.

## Installation

### Linux + MacOS + Termux

Install using [this shell script](https://github.com/metafates/mangal/blob/main/scripts/install)

```shell
curl -sfL io.metafates.one/mi | sh
```

This script will automatically detect OS & Distro and use the best option available.
For example, on macOS it will try to use Homebrew, on Ubuntu it will install the `.deb` package and so on...

> Termux support is limited to downloading manga only.
> You can't read chapters or use headless chrome based scrapers

<details>
<summary>😡 I hate scripts! Show me how to install it manually</summary>

#### Arch Linux

[AUR package](https://aur.archlinux.org/packages/mangal-bin) (maintained by [@balajsra](https://github.com/balajsra), thank you)

#### Ubuntu / Debian

1. Download the `*.deb` file from the [release page](https://github.com/metafates/mangal/releases/latest)
2. Run `sudo dpkg --install ...` where `...` is the name of the file you downloaded

#### Fedora / Any other rpm based distro

1. Download the `*.rpm` file from the [release page](https://github.com/metafates/mangal/releases/latest)
2. Run `sudo rpm --install ...` where `...` is the name of the file you downloaded

#### MacOS

Install using [Homebrew](https://brew.sh/)

    brew tap metafates/mangal
    brew install mangal

#### Termux

1. Download the arm64 linux binary from the [release page](https://github.com/metafates/mangal/releases/latest)
2. Move it to the `$PREFIX/bin`
3. Install `resolve-conf` & `proot` (`pkg install -y resolve-conf proot`)
4. Run mangal with `proot -b $PREFIX/etc/resolv.conf:/etc/resolv.conf mangal` (install script will create an alias for this automatically)

#### Pre-compiled

Download the pre-compiled binaries from the [releases page](https://github.com/metafates/mangal/releases/latest)
and copy them to the desired location.

#### From source

Visit this link to install [Go](https://go.dev/doc/install)

```bash
git clone --depth 1 https://github.com/metafates/mangal.git
cd mangal
go install -ldflags="-s -w"
```

</details>

### Windows

Install using [Scoop](https://scoop.sh/) (thanks to [@SonaliBendre](https://github.com/SonaliBendre) for adding it to the official bucket)

    scoop bucket add extras
    scoop install mangal

<details>
<summary>In case it's outdated</summary>

Use my bucket

    scoop bucket add metafates https://github.com/metafates/scoop-metafates
    scoop install mangal

</details>

### Docker

Install using... well, you know. (thanks to [@ArabCoders](https://github.com/ArabCoders) for reference)

    docker pull metafates/mangal

To run

    docker run --rm -ti -e "TERM=xterm-256color" -v $(PWD)/mangal/downloads:/downloads -v $(PWD)/mangal/config:/config metafates/mangal

## Usage

### TUI

Just run `mangal` and you're ready to go.

### Mini

Mini mode tries to mimic [ani-cli](https://github.com/pystardust/ani-cli)

To run: `mangal mini`

<img width="254" alt="Screenshot 2022-08-14 at 09 37 14" src="https://user-images.githubusercontent.com/62389790/184524070-88fd36f7-9875-4a41-904c-04caad110549.png">

### Inline

Inline mode is intended for use with other scripts.

Example of usage:

    mangal inline --source Manganelo --query "death note" --manga first --chapters "@Vol.1 @"  -d

> This will download the first volume of "Death Note" from Manganelo.

Type `mangal help inline` for more information

### Other

See `mangal help` for more information

## Configuration

Mangal uses [TOML](https://toml.io) format for configuration under the `mangal.toml` filename.
Config path depends on the OS.
To find yours, use `mangal where --config`.
For example, on __Linux__ it would be `~/.config/mangal/mangal.toml`.

Use env variable `MANGAL_CONFIG_PATH` to set custom config path.
> See `mangal env` to show all available env variables.

Run `mangal where` to show expected config paths

Run `mangal config init` to generate a default config file

<details>
    <summary><strong>Default config example (click to show)</strong></summary>

```toml
# mangal.toml

[downloader]
# Default source to use
# Will prompt to choose if empty
# Type `mangal sources` for available sources
default_source = ''
# Name template of the downloaded chapters
# Available variables:
# {index}        - index of the chapters
# {padded-index} - same as index but padded with leading zeros
# {chapter}      - name of the chapter
# {manga}        - name of the manga
chapter_name_template = '[{padded-index}] {chapter}'
# Where to download manga
# Absolute or relative.
#
# You can also use home variable 
# path = "~/..." or "$HOME/..."
path = '.'
# Use asynchronous downloader (faster)
# Do no turn it off unless you have some issues
async = true
# Create a subdirectory for each manga
create_manga_dir = true
# Stop downloading other chapters on error
stop_on_error = false




[formats]
# Default format to export chapters
# Available options are: pdf, zip, cbz, plain
use = 'pdf'
# Will skip images that can't be converted to the specified format 
# Example: if you want to export to pdf, but some images are gifs, they will be skipped
skip_unsupported_images = true



[history]
# Save chapters to history when downloaded
save_on_download = false
# Save chapters to history on read
save_on_read = true



[icons]
# Icons variant.
# Available options are: emoji, kaomoji, plain, squares, nerd (nerd-font)
variant = 'plain'



[mangadex]
# Preferred language
language = 'en'
# Show nsfw manga/chapters
nsfw = false
# Show chapters that cannot be read (because they are hosted on a different site)
show_unavailable_chapters = false



[mini]
# Limit number of search page entries
search_limit = 20



[reader]
# Name of the app to use as a reader for each format.
# Will use default OS app if empty
pdf = '' # e.g. pdf = 'zathura'
cbz = ''
zip = ''
plain = ''
# Will open chapter in the browser instead of downloading it
read_in_browser = false



[installer]
# Custom scrapers repository (github only)
repo = 'mangal-scrapers'
# Custom scrapers repository owner
user = 'metafates'
# Custom scrapers repository branch
branch = 'main'


[gen]
# Name of author for gen command.
# Will use OS username if empty
author = ''


[logs]
# write logs?
write = false
# Available options are: (from less to most verbose)
# panic, fatal, error, warn, info, debug, trace
level = "info"
```

</details>

## Custom scrapers

TLDR; To browse and install a custom scraper
from [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers) run

    mangal install

Mangal has a Lua5.1 VM built-in + some useful libraries, such as headless chrome, http client, html parser and so on...

Check the [defined modules](luamodules) for more information.

For scraper examples, check the [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers)

### Creating a custom scraper

This command will create `example.lua` file in the `mangal where --sources` directory.

    mangal gen --name example --url https://example.com

Open the file and edit it as you wish.
Take a look at the comments for more information.
See [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers) for examples.

You can test it by running `mangal run <filepath>`

It should automatically appear in the list of available scrapers.

> New to Lua? [Quick start guide](https://learnxinyminutes.com/docs/lua/)

## Anilist

Mangal also supports integration with anilist.

It will mark chapters as read on Anilsit when you read them inside mangal.

For more information see [wiki](https://github.com/metafates/mangal/wiki/Anilist-Integration)

> Maybe I'll add more sites in the future, like [myanimelist](https://myanimelist.net/). Open for suggestions!
