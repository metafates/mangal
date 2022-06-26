package main

import "time"

// Version is the current version of the program.
const Version = "1.6.0"

// Mangal is a name of application
// I keep it in a constant to avoid typos
const Mangal = "Mangal"

// CachePrefix is prefix of cache files
const CachePrefix = Mangal + "Cache"

// TempPrefix is prefix of temp files
const TempPrefix = Mangal + "Temp"

// Parallelism is number of parallel workers for scraping
const Parallelism = 100

// TestQuery is a default query for testing
const TestQuery = "Death Note"

// Forever is a constant for inifite time duration.
// It approximates to 292 years
const Forever = time.Duration(1<<63 - 1)

var AvailableFormats = []FormatType{PDF, CBZ, Plain, Zip, Epub}
var FormatsInfo = map[FormatType]string{
	PDF:   "Chapters as PDF with images",
	CBZ:   "Comic book archive format. Basically zip but with .cbz extension",
	Plain: "Just folders with raw .jpg images as chapters",
	Zip:   "Chapters compressed in zip archives",
	Epub:  "eBook format. Packs multiple chapters into single file",
}

// AsciiArt of the app
// I think it looks cool :)
const AsciiArt = "                                _\n" +
	"  /\\/\\   __ _ _ __   __ _  __ _| |\n" +
	" /    \\ / _` | '_ \\ / _` |/ _` | |\n" +
	"/ /\\/\\ \\ (_| | | | | (_| | (_| | |\n" +
	"\\/    \\/\\__,_|_| |_|\\__, |\\__,_|_|\n" +
	"                    |___/         "

// Will be set during build
var (
	AnilistClientSecret string
	AnilistClientID     string
)
