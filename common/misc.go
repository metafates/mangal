package common

import (
	"time"
)

const (

	// Mangal is a name of application
	// I keep it in a constant to avoid typos
	Mangal = "Mangal"

	// CachePrefix is prefix of cache files
	CachePrefix = Mangal + "Cache"

	// TempPrefix is prefix of temp files
	TempPrefix = Mangal + "Temp"

	// Parallelism is number of parallel workers for scraping
	Parallelism = 100

	// TestQuery is a default query for testing
	TestQuery = "Death Note"

	// Forever is a constant for inifite time duration.
	// It approximates to 292 years
	Forever = time.Duration(1<<63 - 1)

	// Referer is a default referer for requests
	Referer = "https://www.google.com"

	// AsciiArt of the app
	// I think it looks cool :)
	AsciiArt = "                                _\n" +
		"  /\\/\\   __ _ _ __   __ _  __ _| |\n" +
		" /    \\ / _` | '_ \\ / _` |/ _` | |\n" +
		"/ /\\/\\ \\ (_| | | | | (_| | (_| | |\n" +
		"\\/    \\/\\__,_|_| |_|\\__, |\\__,_|_|\n" +
		"                    |___/         "
)
