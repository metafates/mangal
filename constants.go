package main

import "time"

const AppName = "Mangal"
const CachePrefix = AppName + "Cache"
const TempPrefix = AppName + "Temp"
const Parallelism = 100
const TestQuery = "Death Note"
const Forever = time.Duration(1<<63 - 1) // 292 years

var AvailableFormats = []FormatType{PDF, CBZ, Plain, Zip, Epub}
var FormatsInfo = map[FormatType]string{
	PDF:   "Chapters as PDF with images",
	CBZ:   "Comic book archive format. Basically zip but with .cbz extension",
	Plain: "Just folders with raw .jpg images as chapters",
	Zip:   "Chapters compressed in zip archives",
	Epub:  "eBook format. Packs multiple chapters into single file",
}
