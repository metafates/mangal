package config

import (
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/color"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/style"
	"github.com/spf13/viper"
)

// Field represents a single config field
type Field struct {
	// Key is the key of the field
	Key string
	// Value is the default value of the field
	Value any
	// Description is the description of the field
	Description string
}

// typeName returns the type of the field without reflection
func (f *Field) typeName() string {
	switch f.Value.(type) {
	case string:
		return "string"
	case int:
		return "int"
	case bool:
		return "bool"
	case []string:
		return "[]string"
	case []int:
		return "[]int"
	default:
		return "unknown"
	}
}

func (f *Field) MarshalJSON() ([]byte, error) {
	field := struct {
		Key         string `json:"key"`
		Value       any    `json:"value"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}{
		Key:         f.Key,
		Value:       f.Value,
		Description: f.Description,
		Type:        f.typeName(),
	}

	return json.Marshal(field)
}

// Pretty format field as string for further cli output
func (f *Field) Pretty() string {
	return fmt.Sprintf(
		`%s
%s: %s = %s
`,
		style.Faint(f.Description),
		style.Fg(color.Purple)(f.Key),
		style.Fg(color.Yellow)(f.typeName()),
		style.Fg(color.Cyan)(fmt.Sprintf("%v", viper.Get(f.Key))),
	)
}

// defaults contains all default values for the config.
// It must contain all fields defined in the constant package.
var defaults = [constant.DefinedFieldsCount]Field{
	{
		constant.DownloaderPath,
		".",
		`Where to download manga
Absolute or relative.
You can also use tilde (~) to refer to your home directory or use env variables.
Examples: ~/... or $HOME/... or ${MANGA_PATH}-mangal`,
	},
	{
		constant.DownloaderChapterNameTemplate,
		"[{padded-index}] {chapter}",
		`Key template of the downloaded chapters
Path forbidden symbols will be replaced with "_"
Available variables:
{index}          - index of the chapters
{padded-index}   - same as index but padded with leading zeros
{chapters-count} - total number of chapters
{chapter}        - name of the chapter
{manga}          - name of the manga
{volume}         - volume of the chapter
{source}         - name of the source`,
	},
	{
		constant.DownloaderAsync,
		true,
		`Use asynchronous downloader (faster)
Do no turn it off unless you have some issues`,
	},
	{
		constant.DownloaderCreateMangaDir,
		true,
		`Create a subdirectory for each manga`,
	},
	{
		constant.DownloaderCreateVolumeDir,
		false,
		`Create a subdirectory for each volume`,
	},
	{
		constant.DownloaderReadDownloaded,
		true,
		"If chapter is already downloaded, read it instead of downloading it to temp",
	},
	{
		constant.DownloaderRedownloadExisting,
		false,
		`Redownload chapters that already exist`,
	},
	{
		constant.DownloaderDefaultSources,
		[]string{},
		`Default sources to use.
Will prompt if not set.
Type "mangal sources list" to show available sources`,
	},
	{
		constant.DownloaderStopOnError,
		false,
		`Stop downloading other chapters on error`,
	},
	{
		constant.DownloaderDownloadCover,
		true,
		`Whether to download manga cover or not`,
	},
	{
		constant.FormatsUse,
		"pdf",
		`Default format to export chapters
Available options are: pdf, zip, cbz, plain`,
	},
	{
		constant.FormatsSkipUnsupportedImages,
		true,
		`Will skip images that can't be converted to the specified format 
Example: if you want to export to pdf, but some images are gifs, they will be skipped`,
	},

	{
		constant.MetadataFetchAnilist,
		true,
		`Fetch metadata from Anilist
It will also cache the results to not spam the API`,
	},

	{
		constant.MetadataComicInfoXML,
		true,
		`Generate ComicInfo.xml file for each chapter`,
	},
	{
		constant.MetadataComicInfoXMLAddDate,
		true,
		`Add series release date to each chapter in ComicInfo.xml file`,
	},
	{
		constant.MetadataComicInfoXMLAlternativeDate,
		false,
		"Use download date instead of series release date in ComicInfo.xml file",
	},
	{
		constant.MetadataComicInfoXMLTagRelevanceThreshold,
		60,
		"Minimum relevance of a tag to be added to ComicInfo.xml file. From 0 to 100",
	},
	{
		constant.MetadataSeriesJSON,
		true,
		`Generate series.json file for each manga`,
	},
	{
		constant.MiniSearchLimit,
		20,
		`Limit of search results to show`,
	},
	{
		constant.IconsVariant,
		"plain",
		`Icons variant.
Available options are: emoji, kaomoji, plain, squares, nerd (nerd-font required)`,
	},
	{
		constant.ReaderPDF,
		"",
		"What app to use to open pdf files",
	},
	{
		constant.ReaderCBZ,
		"",
		"What app to use to open cbz files",
	},
	{
		constant.ReaderZIP,
		"",
		"What app to use to open zip files",
	},
	{
		constant.RaderPlain,
		"",
		"What app to use to open folders",
	},
	{
		constant.ReaderBrowser,
		"",
		"What browser to use to open webpages",
	},
	{
		constant.ReaderFolder,
		"",
		"What app to use to open folders",
	},
	{
		constant.ReaderReadInBrowser,
		false,
		"Open chapter url in browser instead of downloading it",
	},
	{
		constant.HistorySaveOnRead,
		true,
		"Save history on chapter read",
	},
	{
		constant.HistorySaveOnDownload,
		false,
		"Save history on chapter download",
	},
	{
		constant.SearchShowQuerySuggestions,
		true,
		"Show query suggestions in when searching",
	},
	{
		constant.MangadexLanguage,
		"en",
		`Preferred language for mangadex
Use "any" to show all languages`,
	},
	{
		constant.MangadexNSFW,
		false,
		"Show NSFW content",
	},
	{
		constant.MangadexShowUnavailableChapters,
		false,
		"Show chapters that cannot be downloaded",
	},
	{
		constant.InstallerUser,
		"metafates",
		"Custom scrapers repository owner",
	},
	{
		constant.InstallerRepo,
		"mangal-scrapers",
		"Custom scrapers repository name",
	},
	{
		constant.InstallerBranch,
		"main",
		"Custom scrapers repository branch",
	},
	{
		constant.GenAuthor,
		"",
		"Key to use in generated scrapers as author",
	},
	{
		constant.LogsWrite,
		false,
		"Write logs",
	},
	{
		constant.LogsLevel,
		"info",
		`Available options are: (from less to most verbose)
panic, fatal, error, warn, info, debug, trace`,
	},
	{
		constant.LogsJson,
		false,
		"Use json format for logs",
	},
	{
		constant.AnilistEnable,
		false,
		"Enable Anilist integration",
	},
	{
		constant.AnilistCode,
		"",
		"Anilist code to use for authentication",
	},
	{
		constant.AnilistID,
		"",
		"Anilist ID to use for authentication",
	},
	{
		constant.AnilistSecret,
		"",
		"Anilist secret to use for authentication",
	},
	{
		constant.AnilistLinkOnMangaSelect,
		true,
		"Show link to Anilist on manga select",
	},
	{
		constant.TUIItemSpacing,
		1,
		"Spacing between items in the TUI",
	},
	{
		constant.TUIReadOnEnter,
		true,
		"Read chapter on enter if other chapters aren't selected",
	},
	{
		constant.TUISearchPromptString,
		"> ",
		"Search prompt string to use",
	},
	{
		constant.TUIShowURLs,
		true,
		"Show URLs under list items",
	},
	{
		constant.TUIReverseChapters,
		false,
		"Reverse chapters order",
	},
	{
		constant.TUIShowDownloadedPath,
		true,
		"Show path where chapters were downloaded",
	},
}

func init() {
	var count int

	for _, field := range defaults {
		if _, ok := Default[field.Key]; ok {
			panic("Duplicate key in defaults: " + field.Key)
		}

		Default[field.Key] = field
		EnvExposed = append(EnvExposed, field.Key)
		count++
	}

	if count != constant.DefinedFieldsCount {
		panic(fmt.Sprintf("Expected %d default values, got %d", constant.DefinedFieldsCount, count))
	}
}

var Default = make(map[string]Field, constant.DefinedFieldsCount)
