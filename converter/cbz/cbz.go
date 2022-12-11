package cbz

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/key"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/spf13/viper"
	"io"
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

	err = SaveTo(chapter, path)
	if err != nil {
		return "", err
	}

	return path, nil
}

func SaveTo(chapter *source.Chapter, to string) error {
	cbzFile, err := filesystem.Api().Create(to)
	if err != nil {
		return err
	}

	defer util.Ignore(cbzFile.Close)

	zipWriter := zip.NewWriter(cbzFile)
	defer util.Ignore(zipWriter.Close)

	for _, page := range chapter.Pages {
		if err = addToZip(zipWriter, page.Contents, page.Filename()); err != nil {
			return err
		}
	}

	if viper.GetBool(key.MetadataComicInfoXML) {
		comicInfo := chapter.ComicInfo()
		marshalled, err := xml.MarshalIndent(comicInfo, "", "  ")
		if err == nil {
			buf := bytes.NewBuffer(marshalled)
			err = addToZip(zipWriter, buf, "ComicInfo.xml")
		}
	}

	return err
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
