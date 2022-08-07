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

## Features

- __LUAAAA SCRAPPEERRRSS!!!__ You can add any source you want by creating your own _(or using someone's else)_ scraper with __Lua 5.1__.
- __Download & Read Manga__ - I mean, it would be strange if you couldn't, right?
- __4 Different export formats__ - PDF, CBZ, ZIP and plain images
- __Fast__ - yes.
- __Monolith__ - ZERO runtime dependencies. Even Lua is built in.

## Installation

### Linux

Download the latest version from [GitHub release page](https://github.com/metafates/mangal/releases/latest)

Or install it using [Go](https://go.dev/doc/install) 

    go install github.com/metafates/mangal@latest

### macOS

Install using [Homebrew](https://brew.sh/)

    brew tap metafates/mangal
    brew install mangal

### Windows

Install using [Scoop](https://scoop.sh/)

    scoop install ...

### Docker

Install using... well, you know.

    docker pull metafates/mangal

To run

    docker run --rm -ti -e "TERM=xterm-256color" -v (PWD)/mangal/downloads:/downloads -v (PWD)/mangal/config:/config metafates/mangal

## Usage

### TUI

Just run `mangal` and you're ready to go.

### Mini

There's also a `mini` mode that kinda resembles [ani-cli](https://github.com/pystardust/ani-cli)

Run `mangal mini`

## Configuration

Mangal uses [TOML](https://toml.io/en/) format for configuration under the `mangal.toml` filename.
Config is expected to be either at the OS default config directory or under the home directory.
For example, on Linux it would be `~/.config/mangal/mangal.toml` or `~/mangal.toml`.

Run `mangal where` to show expected config paths

> "But what if I want to specify my own config path?"
> 
> Okay, fine, use env variable `MANGAL_CONFIG_PATH`

## Custom scrapers

This is where it gets interesting 😈

Mangal has a Lua5.1 VM built-in + some useful libraries, such as headless chrome, http client, html parser and so on...

Check the [defined modules](luamodules) for more information.
For scraper examples, check the [examples](examples) folder. (feel free to contribute!)

_Okay, so, how do I add a custom scraper?_

1. Create a new lua file in the `mangal where --sources` folder
2. Filename will be used as a source name
3. Your script should contain __3 essential functions__
   - `SearchManga(query)` - must return a table of tables each having 2 fields `name` and `url`
   - `MangaChapters(mangalUrl)` - must return a table of tables each having 2 fields `name` and `url` _(again)_
   - `ChapterPages(chapterUrl)` - must return a table of tables each having 2 fields `index` _(for ordering)_ and `url` _(to download image)_
4. __That's it!__ You can test it by running `mangal run ...` where `...` is a filename

New to Lua? [Quick start guide](https://learnxinyminutes.com/docs/lua/)