package zip

import (
	"archive/zip"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"io"
	"time"
)

type ZIP struct{}

func New() *ZIP {
	return &ZIP{}
}

func (*ZIP) Save(chapter *source.Chapter) (string, error) {
	return save(chapter, false)
}

func (*ZIP) SaveTemp(chapter *source.Chapter) (string, error) {
	return save(chapter, true)
}

func save(chapter *source.Chapter, temp bool) (path string, err error) {
	path, err = chapter.Path(temp)
	if err != nil {
		return
	}

	zipFile, err := filesystem.Api().Create(path)
	if err != nil {
		return
	}

	defer util.Ignore(zipFile.Close)

	zipWriter := zip.NewWriter(zipFile)
	defer util.Ignore(zipWriter.Close)

	for _, page := range chapter.Pages {
		if err = addToZip(zipWriter, page.Contents, page.Filename()); err != nil {
			return "", err
		}
	}

	return
}

func addToZip(writer *zip.Writer, file io.Reader, name string) error {
	header := &zip.FileHeader{
		Name:     name,
		Method:   zip.Deflate,
		Modified: time.Now(),
	}

	headerWriter, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(headerWriter, file)
	return err
}
