# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com), and this project adheres to
[Semantic Versioning](https://semver.org).

## 4.0.6

- Update dependencies
- Fix lua library

## 4.0.5

- Fixes runtime crash #135
- Option to disable colors in cli help `mangal config info -k cli.colored` (why not? =P)
- Improved `config info` command output. It now shows default value and env variable name.
- Internal improvements

## 4.0.4

- Fix manga tags and genres being the same inside ComicInfo.xml #133
- Fill `DescriptionText` field for *series.json* 

## 4.0.3

- Add `exact` manga selector for inline mode #131
- Fix panic when manga not found in inline mode
- More consistent JSON output for inline mode

## 4.0.2

- Fix invalid title in ComicInfo for chapters #130

## 4.0.1

This update includes just some bug-fixes and internal improvements. 🥱 

- Better caching by [Gache](https://github.com/metafates/gache) library
- Fix comic_info_xml_add_date and comic_info_xml_alternative_date flags #126
- Fix notification that new version is available even though it's the same #125
- Fix config set command doesn't work for string values #127
- Fix json output for `config info -j` command #129
- Fix history and default sources not working well together in TUI
- `config reset` now accepts `--all` flag to reset all config values

## 4.0.0

I've been actively working on this update lately, and I'm finally happy to share the 4th version of Mangal! 🐳

The most important feature this major version brings is significantly improved caching mechanism
which makes Mangal extremely fast and responsive.

Now, mangal makes almost no requests to the servers.
This includes Anilist, Scrapers, Update checker and so on!

<details>
<summary><strong>⚠️  BREAKING!!! ⚠️ </strong> Please, read!</summary>

1. `mangal sources` will no longer list available sources, use `mangal sources list` instead.
2. `mangal gen` and `mangal install` were removed. Use `mangal sources gen` and `mangal sources install` instead.
3. `mangal sources remove` command improved and accepts flags instead of args.

Inline JSON output is different now.

- JSON fields now follow the [camelCase](https://en.wikipedia.org/wiki/Camel_case) style instead of `PascalCase`
  (actually, using PascalCase was never a goal, I just forgot to properly configure it).
  But since it's a major release I can finally fix this.
- Structure was changed
- Additional fields were added

See [Inline mode wiki](https://github.com/metafates/mangal/wiki/Inline-mode) for new output schemas.

Please, consider these changes when migrating your applications that use mangal from 3rd version to 4th.
</details>

- Improved TUI experience
- Search completions in TUI. `mangal config info -k search.show_query_suggestions`
- Anilist caching significantly improved. Now, it will cache all search results (for 2 days)
- Update metadata of already downloaded manga (ComicInfo.xml, series.json, cover image) after changing Anilist bind. #124
  See `mangal inline anilist update` for more info
- New command to generate json schema of inline output. See `mangal help inline schema`
- **Breaking** `downloader.default_source` was changed to `downloader.default_sources` and accepts array of strings.
  See `mangal config info -k downloader.default_sources` for more info
- New `config reset` command
- Add caching for custom (lua) sources
- Include different cover sizes and color for json output #116
- Add option to omit dates for ComicInfo.xml #117
- By default, when reading a chapter, mangal will look for its downloaded copy, instead of downloading it again.
  See `mangal config info -k downloader.read_downloaded`
- Overwrite old `series.json` file each time a chapter is downloaded
- Detect sources that use headless chrome and show that in the item description when selecting sources
- Option to use alternative ComicInfo.xml date.
  See `mangal config info -k metadata.comic_info_xml_alternative_date` for more info
- Notify about new version in `help` command
- Include staff in ComicInfo.xml #119
- Add `--set-only` and `--unset-only` flags for `env` command. Old `--filter` flag was removed
- `version` command now has `--short` to just print the version without extra information
- **Breaking!** Your old reading history (via `mangal --continue`) will be reset
- Improved `clear` command
- Option to set threshold for tag relevance to be included in ComicInfo.xml #121
- Improved inline command json output, fixes
- Internal improvements

Enjoy!

## 3.14.2

- Do not put an invalid value for dates #114
- Set `metadata.series_json` to `true` by default.
  See `mangal config info -k metadata.series_json` for more info

## 3.14.1

- Mark flags as required for `inline anilist` commands
- Remove `update` command [Why?](https://github.com/metafates/mangal/discussions/112)
- `mangal version` will notify if new version is available
- Use correct page image extension for custom sources #110

## 3.14.0

- New commands related to the anilist manga linkage. Now you can set what anilist manga should be linked with what titles by id. See `mangal inline anilist` for more information. #106
- Increase default http timeout to 20 seconds #108
- Fixed nil panic when trying to resume reading from history with mini mode

## 3.13.0

- Support environment variables for `downloader.path` config field #103
- Replace Mangakakalot with Manganato #102
- Move `install` & `gen` commands to `sources` subcommands. E.g. if you used `mangal install` before use `mangal sources install`. Old commands are still present, but marked as deprecated.
- New flags `--builtin` & `--custom` for `sources` command to filter sources by type.
- New flag `--json` added for `config info` command to show fields in json format. 
- New command `mangal sources remove <name>` to remove custom source.
- Minor performance improvements.

## 3.12.0

- Faster and more optimized page downloader
- Show current config field value in `config info` cmd
- Optimize PDF converter
- By default, mangal will not redownload chapters that already exists at the target path. Can be disabled with `mangal config set -k downloader.redownload_existing -v true` #100 
- Sort config fields in `config info` command
- Better looking `version` command

## 3.11.1

- Fixed critical bug when mangal would crash when using mini mode
- Slightly change `version` command output
- Use better text wrapping

## 3.11.0

- Add an option to search mangas with inline mode. `mangal inline -q "..." -j` will output search results without chapters. #97
- `config` cmd improved. Now, `config set` will automatically parse the value to  the expected type. 
- Internal improvements.


## 3.10.0

- New feature: you can choose what anilist manga to link by pressing <kbd>a</kbd> in the manga chapters list.
  TUI mode only. This will affect what metadata is downloaded for the manga and what manga would be marked as read on your anilist profile.
- Pressing <kbd>enter</kbd> on the chapter will now open it for reading if other chapters aren't selected.
  Can be disabled with `mangal config set -k tui.read_on_enter -bv false`
- The chapter selection page now shows which manga from the anilist it is linked to.
  Can be disabled with `mangal config set -k anilist.link_on_manga_select -bv false`
- Add an option to change spacing between items in the TUI.
  Can be changed with `mangal config set -k tui.item_spacing -iv 1` (1 is default)
- List filtering in the TUI works better now by stripping the icon
- Option to hide list items urls in TUI. To hide: `mangal config set -k tui.show_urls -bv false`
- After downloading of chapters is done, mangal will show the output path.
  To disable: `mangal config set -k tui.show_downloaded_path -bv false`
- Option to reverse order of the chapters in the TUI. `mangal config set -k tui.reverse_chapters -bv true` to enable
- Reduce the size of the compiled binary by removing unused lua libraries. May break some lua scripts, but I don't think you were using AWS to scrape manga :)

## 3.9.1

- Fix version comparison mechanism for `update` command.
  Now it compares each fragment separately (major, minor, patch) instead of comparing two versions as strings lexicographically.

## 3.9.0

- New sources: [Mangakakalot](https://mangakakalot.com) & [Mangapill](https://mangapill.com)
- Fix termux installation detection by [@2096779623](https://github.com/2096779623) #94
- Change the way `mangal update` works.
  If mangal wasn't installed via package manager it will get the current location of running binary and replace it with the new one.
  Previously it was assumed that mangal was installed in `/usr/local/bin` which is not always the case.
  Doesn't work on Termux yet because it requires specific variation of Go compiler which is troublesome to configure with automatic release system that mangal uses.
  Will be fixed in the future.

## 3.8.1

- Fix installation method detection
- Fix install script

## 3.8.0

- Support for more manga metadata fields such as summary, genres, tags, and more.
- Fetch manga metadata from anilist.
  `metadata.fetch_anilist` (default: `true`) 
- Generate `series.json` file.
  `metadata.series_json` (default: `false`)
- Generate `ComicInfo.xml` file (for CBZ only)
  `metadata.comic_info_xml` (default: `true`)
- Support for downloading manga covers.
  `downloader.download_cover` (default: `false`)
- Better progress message while downloading in TUI mode
- Set option `downloader.create_volume_dir` to `false` by default
- Version command now shows more information (such as build date, commit hash, etc.)
- New flag for inline mode: `--output/-o` to redirect output to file (will use STDOUT if not set)
- New `mangal config set` command to set config values. See `mangal help config set` for more info.
- New `mangal config get` command to get config values. See `mangal help config get` for more info 
- New `mangal config info` command to list all available config fields with description for each.
- Improve `mangal clear` command. It's more accurate and faster now. 
- Better cache & temp files handling
- Fix `mangal update` command when it was not able to update using script. 
- Expose every config field as ENV variable (see `mangal env` to show all of them)

## 3.7.0

- Add support for volumes - now you can select chapters by volume (in TUI mode only). #90
- New config field `downloader.create_volume_dir` to create a subdirectory for each volume if it's known
  (default: `true`)
- New feature - search with multiple sources at once (TUI mode only). #86
- New config field `logs.json` to write logs in json format (default: `false`)
- Better keymap help
- Slightly more logs
- Fix sources command: do not print `custom sources` header if there are none
- More minor fixes

## 3.6.0

- Add `--json` flag for the `inline` mode by [@jojoxd](https://github.com/jojoxd) #83
- Fixed `mangal update` command
- Rename `history` global flag to `write-history` (shorthands are the same `-H`)
- Add basic Termux support. Scripts that use headless Chrome browser won't work as well as reading mode. That means that
  you can use it only for downloading manga (for now at least) #80
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
- Fixed bug where chapters with colon in title would not open for read on Windows #24
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
