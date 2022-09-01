# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com), and this project adheres to
[Semantic Versioning](https://semver.org).

## 3.6.0

- Add `--json` flag for the `inline` mode by [@jojoxd](https://github.com/jojoxd) #83
- Fixed `mangal update` command
- Rename `history` global flag to `write-history` (shorthands are the same `-H`)
- Add basic Termux support. Scripts that use headless Chrome browser won't work as well as reading mode. That means that you can use it only for downloading manga (for now at least) #80
- Fixed a bug where reading manga with mini mode would cause mangal to crash #82


## 3.5.0

- `mangal update` command added to update itself
- `mangal gen` command added to generate a template for custom source
- Added `--raw` flag for the `sources` command to print without headers
- Added `--downloads` flag for the `where` command to print location of the downloads directory
- Option to show all languages in mangadex by setting `mangadex.language = 'any'` in the config.
- Show chapter index in the TUI history list
- Fixed bug where pressing <kbd>confirm</kbd> button with empty history would cause mangal to crash
- Fixed bug where error message would not fit the screen
- Fixed bug when `mangal config init --force` would ignore `--force` flag
- Internal performance improvements

## 3.4.1

- Option to continue downloading chapters after fail
- Option to redownload failed chapters
- Option to select custom reader app for each format
- Option to skip images with unsupported formats by converter (e.g. pdf converter will skip .gif images) (#77)
- Option to specify custom repository for `mangal install` command
- Fixed error when using custom readers was not possible
- Highlight <kbd>read</kbd> button in the chapters list
- Make `mangal env` more compact, add `--filter` argument to filter out unset variables
- Show `MANGAL_CONFIG_PATH` in the `env` command
- Show better error when chapter could not be opened by reader
- Fix chapters range selector in inline mode (#76)
- Show progress when initializing source from history
- Show size of a downloading chapter

## 3.3.0

- Inline mode added. See `mangal help inline`
- Option to choose default source added
- Show `read` button in help menu
- Bug fixes and improvements

## 3.2.1

- Fix home variable in config

## 3.2.0

- New command added mangal install to install scrapers from mangal-scrapers repo
- Added an option to remove a single entry from history by pressing d
- Added an option to download chapters without creating manga directory
- Dependencies updated
- Bug fixes and improvements

## 3.1.0

- `where` command now prints to stdout. It can be used like that: `cd $(mangal where --config)`
- Mini mode was completely rewritten to look exactly like [ani-cli](https://github.com/pystardust/ani-cli)
- PDF is a default export format now (was plain)
- Plain icons are default now (were emoji)
- New icons added - "squares"
- New command `mangal config remove` to... well... remove config
- Minor bug fixes and improvements

## 3.0.3

- Better path handling
- Use pdf a default format

## 3.0.2

- Fix bug where empty config would case errors

## 3.0.1

- Bug fixes...

## 3.0.0

- Full rewrite of the mangal
- Add support for Lua scrapers
- Better TUI
- Mini mode

## 2.2.0

- History mode added. Now you can resume your readings by launching mangal with `mangal --resume` flag (`-r` for short)
- Support for new environment variables added
- `mangal env` command added to list all env variables and their values (if any set)
- ComicInfo.xml file fixed #53

## 2.1.1

- `doctor` command now shows more information about errors
- minor bug fixes and performance improvements

## 2.1.0

- Significant performance improvements! 🚀
- Reduced disk usage
- Add support for env variables `MANGAL_CONFIG_PATH` and `MANGAL_DOWNLOAD_PATH`
- Improved config structure (breaking change, reinitialize your config if you have one)
- ComicInfo.xml support added for CBZ format #27
- `config init --clean` flag added that creates config without additional comments
- `config remove` command added to delete user config

## 2.0.1

- Fixed #36
- Small shell completion improvements

## 2.0.0

- Anilist integration BETA
- Diagnostics command `mangal doctor`
- `mangal check-update` renamed to `mangal latest`
- Custom naming templates for chapters (like this `[%d] %s`)
- Bug fixes
- Faster config parser (up to 5x faster!)
- Minor improvements

## 1.5.2

- Command to check for update added `mangal check-update` #26
- Scraper system improved
- Fixed bug where chapters with colon in title would not open for read on windows #24
- Various bug fixes

## 1.5.1

- Fixes #21, #20

## 1.5.0

- Epub format added
- Move UI related configurations to [ui] section in the config file
- New command `formats` added to show available formats and their description
- Minor improvements

## 1.4.2

- Fixes #15

## 1.4.1

- Multiple formats support
