package main

import "time"

const AppName = "Mangal"

// CachePrefix is prefix of cache files
const CachePrefix = AppName + "Cache"

// TempPrefix is prefix of temp files
const TempPrefix = AppName + "Temp"

// Parallelism is number of parallel workers for scraping
const Parallelism = 100

// TestQuery is a default query for testing
const TestQuery = "Death Note"

// Forever is a constant for inifite time duration
const Forever = time.Duration(1<<63 - 1) // 292 years

var AvailableFormats = []FormatType{PDF, CBZ, Plain, Zip, Epub}
var FormatsInfo = map[FormatType]string{
	PDF:   "Chapters as PDF with images",
	CBZ:   "Comic book archive format. Basically zip but with .cbz extension",
	Plain: "Just folders with raw .jpg images as chapters",
	Zip:   "Chapters compressed in zip archives",
	Epub:  "eBook format. Packs multiple chapters into single file",
}
