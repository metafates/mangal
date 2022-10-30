package source

import "encoding/xml"

type ComicInfo struct {
	XMLName  xml.Name `xml:"ComicInfo"`
	XmlnsXsi string   `xml:"xmlns:xsi,attr"`
	XmlnsXsd string   `xml:"xmlns:xsd,attr"`

	// General
	Title      string `xml:"Title"`
	Series     string `xml:"Series"`
	Number     int    `xml:"Number"`
	Web        string `xml:"Web"`
	Genre      string `xml:"Genre"`
	PageCount  int    `xml:"PageCount"`
	Summary    string `xml:"Summary"`
	Count      int    `xml:"Count"`
	Characters string `xml:"Characters"`
	Year       int    `xml:"Year,omitempty"`
	Month      int    `xml:"Month,omitempty"`
	Day        int    `xml:"Day,omitempty"`
	Writer     string `xml:"Writer,omitempty"`
	Penciller  string `xml:"Penciller,omitempty"`
	Letterer   string `xml:"Letterer,omitempty"`
	Translator string `xml:"Translator,omitempty"`
	Tags       string `xml:"Tags"`
	Notes      string `xml:"Notes"`
	Manga      string `xml:"Manga"`
}
