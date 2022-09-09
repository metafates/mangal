package cbz

import (
	"archive/zip"
	"bytes"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
	"io"
	"strings"
	"text/template"
)

type CBZ struct{}

func New() *CBZ {
	return &CBZ{}
}

func (*CBZ) Save(chapter *source.Chapter) (string, error) {
	return save(chapter, false)
}

func (*CBZ) SaveTemp(chapter *source.Chapter) (string, error) {
	return save(chapter, true)
}

func save(chapter *source.Chapter, temp bool) (path string, err error) {
	path, err = chapter.Path(temp)
	if err != nil {
		return
	}

	cbzFile, err := filesystem.Get().Create(path)
	if err != nil {
		return
	}

	defer util.Ignore(cbzFile.Close)

	zipWriter := zip.NewWriter(cbzFile)
	defer util.Ignore(zipWriter.Close)

	for _, page := range chapter.Pages {
		if err = addToZip(zipWriter, page.Contents, page.Filename()); err != nil {
			return
		}
	}

	err = addToZip(zipWriter, comicInfo(chapter), "ComicInfo.xml")
	return
}

func comicInfo(chapter *source.Chapter) *bytes.Buffer {
	t := `
<ComicInfo xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<Title>{{ .Name }}</Title>
  	<Series>{{ .Manga.Name }}</Series>
  	<Web>{{ .Manga.URL }}</Web>
	<Genre>{{ join .Manga.Genres "," }}</Genre>
	<PageCount>{{ len .Pages }}</PageCount>
	<Summary>{{ .Manga.Summary }}</Summary>
	<Count>{{ len .Manga.Chapters }}</Count>
	<Writer>{{ .Manga.Author }}</Writer>
  	<Manga>YesAndRightToLeft</Manga>
</ComicInfo>`

	funcs := template.FuncMap{
		"join": strings.Join,
	}

	parsed := lo.Must(template.New("ComicInfo").Funcs(funcs).Parse(t))
	buf := bytes.NewBufferString("")
	lo.Must0(parsed.Execute(buf, chapter))

	return buf
}

func addToZip(writer *zip.Writer, file io.Reader, name string) error {
	header := &zip.FileHeader{
		Name:   name,
		Method: zip.Store,
	}

	headerWriter, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(headerWriter, file)
	return err
}
