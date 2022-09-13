package cbz

import (
	"archive/zip"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
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

	cbzFile, err := filesystem.Api().Create(path)
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

	if viper.GetBool(constant.MetadataComicInfoXML) {
		err = addToZip(zipWriter, chapter.ComicInfoXML(), "ComicInfo.xml")
	}

	return
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
