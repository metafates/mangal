package source

import "encoding/xml"

type ComicInfo struct {
	XMLName  xml.Name `xml:"ComicInfo"`
	XmlnsXsi string   `xml:"xmlns:xsi,attr"`
	XmlnsXsd string   `xml:"xmlns:xsd,attr"`

	// General
	Title      string `xml:"Title,omitempty"`
	Series     string `xml:"Series,omitempty"`
	Number     int    `xml:"Number,omitempty"`
	Web        string `xml:"Web,omitempty"`
	Genre      string `xml:"Genre,omitempty"`
	PageCount  int    `xml:"PageCount,omitempty"`
	Summary    string `xml:"Summary,omitempty"`
	Count      int    `xml:"Count,omitempty"`
	Characters string `xml:"Characters,omitempty"`
	Year       int    `xml:"Year,omitempty"`
	Month      int    `xml:"Month,omitempty"`
	Day        int    `xml:"Day,omitempty"`
	Writer     string `xml:"Writer,omitempty"`
	Penciller  string `xml:"Penciller,omitempty"`
	Letterer   string `xml:"Letterer,omitempty"`
	Translator string `xml:"Translator,omitempty"`
	Tags       string `xml:"Tags,omitempty"`
	Notes      string `xml:"Notes,omitempty"`
	Manga      string `xml:"Manga,omitempty"`
}
