package config

import (
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/color"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/key"
	"github.com/metafates/mangal/style"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"reflect"
	"strconv"
	"strings"
	"text/template"
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
		Default     any    `json:"default"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}{
		Key:         f.Key,
		Value:       viper.Get(f.Key),
		Default:     f.Value,
		Description: f.Description,
		Type:        f.typeName(),
	}

	return json.Marshal(field)
}

var prettyTemplate = lo.Must(template.New("pretty").Funcs(template.FuncMap{
	"faint":  style.Faint,
	"bold":   style.Bold,
	"purple": style.Fg(color.Purple),
	"blue":   style.Fg(color.Blue),
	"cyan":   style.Fg(color.Cyan),
	"value":  func(k string) any { return viper.Get(k) },
	"hl": func(v any) string {
		switch value := v.(type) {
		case bool:
			b := strconv.FormatBool(value)
			if value {
				return style.Fg(color.Green)(b)
			}

			return style.Fg(color.Red)(b)
		case string:
			return style.Fg(color.Yellow)(value)
		default:
			return fmt.Sprint(value)
		}
	},
	"typename": func(v any) string { return reflect.TypeOf(v).String() },
}).Parse(`{{ faint .Description }}
{{ blue "Key:" }}     {{ purple .Key }}
{{ blue "Env:" }}     {{ .Env }}
{{ blue "Value:" }}   {{ hl (value .Key) }}
{{ blue "Default:" }} {{ hl (.Value) }}
{{ blue "Type:" }}    {{ typename .Value }}`))

func (f *Field) Pretty() string {
	var b strings.Builder

	lo.Must0(prettyTemplate.Execute(&b, f))

	return b.String()
}

func (f *Field) Env() string {
	env := strings.ToUpper(EnvKeyReplacer.Replace(f.Key))
	appPrefix := strings.ToUpper(constant.Mangal + "_")

	if strings.HasPrefix(env, appPrefix) {
		return env
	}

	return appPrefix + env
}

// Pretty format field as string for further cli output
//func (f *Field) Pretty() string {
//	return fmt.Sprintf(
//		`%s
//%s: %s = %s
//`,
//		style.Faint(f.Description),
//		style.Fg(color.Purple)(f.Key),
//		style.Fg(color.Yellow)(f.typeName()),
//		style.Fg(color.Cyan)(fmt.Sprintf("%v", viper.Get(f.Key))),
//	)
//}

// defaults contains all default values for the config.
// It must contain all fields defined in the constant package.
var defaults = [key.DefinedFieldsCount]Field{
	{
		key.DownloaderPath,
		".",
		`Where to download manga
Absolute or relative.
You can also use tilde (~) to refer to your home directory or use env variables.
Examples: ~/... or $HOME/... or ${MANGA_PATH}-mangal`,
	},
	{
		key.DownloaderChapterNameTemplate,
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
		key.DownloaderAsync,
		true,
		`Use asynchronous downloader (faster)
Do no turn it off unless you have some issues`,
	},
	{
		key.DownloaderCreateMangaDir,
		true,
		`Create a subdirectory for each manga`,
	},
	{
		key.DownloaderCreateVolumeDir,
		false,
		`Create a subdirectory for each volume`,
	},
	{
		key.DownloaderReadDownloaded,
		true,
		"If chapter is already downloaded, read it instead of downloading it to temp",
	},
	{
		key.DownloaderRedownloadExisting,
		false,
		`Redownload chapters that already exist`,
	},
	{
		key.DownloaderDefaultSources,
		[]string{},
		`Default sources to use.
Will prompt if not set.
Type "mangal sources list" to show available sources`,
	},
	{
		key.DownloaderStopOnError,
		false,
		`Stop downloading other chapters on error`,
	},
	{
		key.DownloaderDownloadCover,
		true,
		`Whether to download manga cover or not`,
	},
	{
		key.FormatsUse,
		"pdf",
		`Default format to export chapters
Available options are: pdf, zip, cbz, plain`,
	},
	{
		key.FormatsSkipUnsupportedImages,
		true,
		`Will skip images that can't be converted to the specified format 
Example: if you want to export to pdf, but some images are gifs, they will be skipped`,
	},

	{
		key.MetadataFetchAnilist,
		true,
		`Fetch metadata from Anilist
It will also cache the results to not spam the API`,
	},

	{
		key.MetadataComicInfoXML,
		true,
		`Generate ComicInfo.xml file for each chapter`,
	},
	{
		key.MetadataComicInfoXMLAddDate,
		true,
		`Add series release date to each chapter in ComicInfo.xml file`,
	},
	{
		key.MetadataComicInfoXMLAlternativeDate,
		false,
		"Use download date instead of series release date in ComicInfo.xml file",
	},
	{
		key.MetadataComicInfoXMLTagRelevanceThreshold,
		60,
		"Minimum relevance of a tag to be added to ComicInfo.xml file. From 0 to 100",
	},
	{
		key.MetadataSeriesJSON,
		true,
		`Generate series.json file for each manga`,
	},
	{
		key.MiniSearchLimit,
		20,
		`Limit of search results to show`,
	},
	{
		key.IconsVariant,
		"plain",
		`Icons variant.
Available options are: emoji, kaomoji, plain, squares, nerd (nerd-font required)`,
	},
	{
		key.ReaderPDF,
		"",
		"What app to use to open pdf files",
	},
	{
		key.ReaderCBZ,
		"",
		"What app to use to open cbz files",
	},
	{
		key.ReaderZIP,
		"",
		"What app to use to open zip files",
	},
	{
		key.RaderPlain,
		"",
		"What app to use to open folders",
	},
	{
		key.ReaderBrowser,
		"",
		"What browser to use to open webpages",
	},
	{
		key.ReaderFolder,
		"",
		"What app to use to open folders",
	},
	{
		key.ReaderReadInBrowser,
		false,
		"Open chapter url in browser instead of downloading it",
	},
	{
		key.HistorySaveOnRead,
		true,
		"Save history on chapter read",
	},
	{
		key.HistorySaveOnDownload,
		false,
		"Save history on chapter download",
	},
	{
		key.SearchShowQuerySuggestions,
		true,
		"Show query suggestions in when searching",
	},
	{
		key.MangadexLanguage,
		"en",
		`Preferred language for mangadex
Use "any" to show all languages`,
	},
	{
		key.MangadexNSFW,
		false,
		"Show NSFW content",
	},
	{
		key.MangadexShowUnavailableChapters,
		false,
		"Show chapters that cannot be downloaded",
	},
	{
		key.InstallerUser,
		"metafates",
		"Custom scrapers repository owner",
	},
	{
		key.InstallerRepo,
		"mangal-scrapers",
		"Custom scrapers repository name",
	},
	{
		key.InstallerBranch,
		"main",
		"Custom scrapers repository branch",
	},
	{
		key.GenAuthor,
		"",
		"Key to use in generated scrapers as author",
	},
	{
		key.LogsWrite,
		false,
		"Write logs",
	},
	{
		key.LogsLevel,
		"info",
		`Available options are: (from less to most verbose)
panic, fatal, error, warn, info, debug, trace`,
	},
	{
		key.LogsJson,
		false,
		"Use json format for logs",
	},
	{
		key.AnilistEnable,
		false,
		"Enable Anilist integration",
	},
	{
		key.AnilistCode,
		"",
		"Anilist code to use for authentication",
	},
	{
		key.AnilistID,
		"",
		"Anilist ID to use for authentication",
	},
	{
		key.AnilistSecret,
		"",
		"Anilist secret to use for authentication",
	},
	{
		key.AnilistLinkOnMangaSelect,
		true,
		"Show link to Anilist on manga select",
	},
	{
		key.TUIItemSpacing,
		1,
		"Spacing between items in the TUI",
	},
	{
		key.TUIReadOnEnter,
		true,
		"Read chapter on enter if other chapters aren't selected",
	},
	{
		key.TUISearchPromptString,
		"> ",
		"Search prompt string to use",
	},
	{
		key.TUIShowURLs,
		true,
		"Show URLs under list items",
	},
	{
		key.TUIReverseChapters,
		false,
		"Reverse chapters order",
	},
	{
		key.TUIShowDownloadedPath,
		true,
		"Show path where chapters were downloaded",
	},
	{
		key.CliColored,
		true,
		"Use colors in CLI help page",
	},
	{
		key.CliVersionCheck,
		true,
		"Check for a new version of the CLI occasionally",
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

	if count != key.DefinedFieldsCount {
		panic(fmt.Sprintf("Expected %d default values, got %d", key.DefinedFieldsCount, count))
	}
}

var Default = make(map[string]Field, key.DefinedFieldsCount)
