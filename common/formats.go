package common

// FormatType is type of format used for output
type FormatType string

const (
	PDF   FormatType = "pdf"
	CBZ   FormatType = "cbz"
	Zip   FormatType = "zip"
	Plain FormatType = "plain"
	Epub  FormatType = "epub"
)

var AvailableFormats = []FormatType{PDF, CBZ, Zip, Plain, Epub}
var FormatsInfo = map[FormatType]string{
	PDF:   "Chapters as PDF with images",
	CBZ:   "Comic book archive format. Basically zip but with .cbz extension",
	Plain: "Just folders with raw .jpg images as chapters",
	Zip:   "Chapters compressed in zip archives",
	Epub:  "eBook format. Packs multiple chapters into single file",
}
