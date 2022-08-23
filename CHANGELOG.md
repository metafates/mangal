# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com), and this project adheres to
[Semantic Versioning](https://semver.org).

## [3.5.0]

- `gen` command added to generate a template for custom source

## [3.4.1]

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
