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

https://user-images.githubusercontent.com/62389790/190131616-09d1f134-7ccd-4f0f-bd4d-15517b8dd2a0.mp4

## Try it!

```shell
curl -sL mangal.metafates.one/run | sh
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
- [Honorable mentions](#honorable-mentions)

## Features

- __Lua Scrapers!!!__ You can add any source you want by creating your own _(or using someone's else)_ scraper with
  __Lua 5.1__. See [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers)
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
curl -sL mangal.metafates.one/install | sh
```

This script will automatically detect OS & Distro and use the best option available.
For example, on macOS it will try to use Homebrew, on Ubuntu it will install the `.deb` package and so on...

> Termux support is limited to downloading manga only.
> You can't read chapters or use headless chrome based scrapers

<details>
<summary>😡 I hate scripts! Show me how to install it manually</summary>

#### Arch Linux

[AUR package](https://aur.archlinux.org/packages/mangal-bin) (maintained by [@balajsra](https://github.com/balajsra),
thank you)

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
4. Run mangal with `proot -b $PREFIX/etc/resolv.conf:/etc/resolv.conf mangal` (install script will create an alias for
   this automatically)

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

Install using [Scoop](https://scoop.sh/) (thanks to [@SonaliBendre](https://github.com/SonaliBendre) for adding it to
the official bucket)

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

```shell
docker run --rm -ti -e "TERM=xterm-256color" -v $(PWD)/mangal/downloads:/downloads -v $(PWD)/mangal/config:/config metafates/mangal
```

## Usage

### TUI

Just run `mangal` and you're ready to go.

<details>
<summary>Keybinds</summary>

| Bind                                                        | Description   |
|-------------------------------------------------------------|---------------|
| <kbd>?</kbd>                                                | Show help     |
| <kbd>↑/j</kbd> <kbd>↓/k</kbd> <kbd>→/l</kbd> <kbd>←/h</kbd> | Navigate      |
| <kbd>g</kbd>                                                | Go to first   |
| <kbd>G</kbd>                                                | Go to last    |
| <kbd>/</kbd>                                                | Filter        |
| <kbd>esc</kbd>                                              | Back          |
| <kbd>space</kbd>                                            | Select one    |
| <kbd>tab</kbd>                                              | Select all    |
| <kbd>v</kbd>                                                | Select volume |
| <kbd>backspace</kbd>                                        | Unselect all  |
| <kbd>enter</kbd>                                            | Confirm       |
| <kbd>o</kbd>                                                | Open URL      |
| <kbd>r</kbd>                                                | Read          |
| <kbd>q</kbd>                                                | Quit          |
| <kbd>ctrl+c</kbd>                                           | Force quit    |

</details>

### Mini

Mini mode tries to mimic [ani-cli](https://github.com/pystardust/ani-cli)

To run: `mangal mini`

<img width="254" alt="Screenshot 2022-08-14 at 09 37 14" src="https://user-images.githubusercontent.com/62389790/184524070-88fd36f7-9875-4a41-904c-04caad110549.png">

### Inline

Inline mode is intended for use with other scripts.

Example of usage:

    mangal inline --source Manganelo --query "death note" --manga first --chapters all  -d

> This will download all chapters of the "Death Note" from Manganelo.

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

| Command              | Description                                      |
|----------------------|--------------------------------------------------|
| `mangal config get`  | Get config value for specific key                |
| `mangal config set`  | Set config value for specific key                |
| `mangal config info` | List all config fields with description for each |
| `mangal config init` | Write current config to a file                   |

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

## Honorable mentions

### Similar Projects

- [mangadesk](https://github.com/darylhjd/mangadesk) - Terminal client for MangaDex
- [ani-cli](https://github.com/pystardust/ani-cli) - A cli tool to browse and play anime
- [manga-py](https://github.com/manga-py/manga-py) - Universal manga downloader
- [animdl](https://github.com/justfoolingaround/animdl) - A highly efficient, fast, powerful and light-weight anime
  downloader and streamer
- [tachiyomi](https://github.com/tachiyomiorg/tachiyomi) - Free and open source manga reader for Android

### Libraries

- [bubbletea](https://github.com/charmbracelet/bubbletea), [bubbles](https://github.com/charmbracelet/bubbles)
  & [lipgloss](https://github.com/charmbracelet/lipgloss) - Made mangal shine! The best TUI libraries ever ✨
- [gopher-lua](https://github.com/yuin/gopher-lua) - Made it possible to write custom scrapers with Lua ❤️
- [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper) - Responsible for the awesome CLI
  & config experience 🛠
- [pdfcpu](https://github.com/pdfcpu/pdfcpu) - Fast pdf processor in pure go 📄
- _And many others!_

### Contributors

And of course, thanks to all the contributors! You are awesome!

<p align="center">
<a href="https://github.com/metafates/mangal/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=metafates/mangal" />
</a>
</p>

---

If you find this project useful or want to say thank you,
please consider starring it, that would mean a lot to me ⭐
