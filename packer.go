package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/bmaupin/go-epub"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/afero"
	"io"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
)

type PackerContext struct {
	Manga   *URL
	Chapter *URL
}

// PackToPDF packs chapter to .pdf format
func PackToPDF(images []*bytes.Buffer, destination string, context *PackerContext) (string, error) {
	destination += ".pdf"

	err := RemoveIfExists(destination)
	if err != nil {
		return "", err
	}

	// Create parent directory since pdfcpu have some troubles when it doesn't exist
	if exists, err := Afero.Exists(filepath.Dir(destination)); err != nil {
		return "", err
	} else if !exists {
		if err := Afero.MkdirAll(filepath.Dir(destination), 0777); err != nil {
			return "", err
		}
	}

	// create destination file
	pdf, err := Afero.Create(destination)

	if err != nil {
		return "", err
	}

	defer func() {
		_ = pdf.Close()
	}()

	// convert []*buffer to []Reader
	var readers = make([]io.Reader, len(images))
	for i, buffer := range images {
		readers[i] = buffer
	}

	err = pdfcpu.ImportImages(nil, pdf, readers, nil, nil)
	if err != nil {
		return "", err
	}

	return destination, nil
}

// EpubFile global epub file that is used to save multiple chapters in a single file
var EpubFile *epub.Epub

// PackToEpub adds chapter to epub file
func PackToEpub(images []*bytes.Buffer, destination string, context *PackerContext) (string, error) {
	var (
		coverSet     = true
		mangaTitle   string
		chapterTitle string
	)

	if context != nil {
		mangaTitle = context.Manga.Info
		chapterTitle = context.Chapter.Info
	} else {
		mangaTitle = filepath.Base(filepath.Dir(destination))
		chapterTitle = filepath.Base(destination)
	}

	destination = filepath.Join(filepath.Dir(destination), mangaTitle+".epub")

	// Initialize epub file
	if EpubFile == nil {
		coverSet = false
		EpubFile = epub.NewEpub(mangaTitle)
		// set reading direction right to left (default for manga)
		EpubFile.SetPpd("rtl")

		// remove epub file if it already exists
		if err := RemoveIfExists(destination); err != nil {
			return "", err
		}
		if err := Afero.MkdirAll(filepath.Dir(destination), 0777); err != nil {
			return "", err
		}

		file, err := Afero.Create(destination)
		if err != nil {
			return "", err
		}

		if err = file.Close(); err != nil {
			return "", err
		}
	}

	var epubImages = make([]string, len(images))

	// add images to epub file
	for i, image := range images {
		// save image to temporary file
		imagePath, err := SaveTemp(image)

		if err != nil {
			return "", err
		}

		epubImage, err := EpubFile.AddImage(imagePath, strconv.Itoa(rand.Intn(100000))+filepath.Base(imagePath))
		if err != nil {
			return "", err
		}
		if !coverSet {
			EpubFile.SetCover(epubImage, "")
			coverSet = true
		}
		epubImages[i] = epubImage
	}

	imgTags := Map(epubImages, func(pathToImage string) string {
		return fmt.Sprintf(`
			<p style="display:block;margin:0;">
				<img src="%s" style="height:auto;width:auto;"/>
            </p>
		`, pathToImage)
	})

	_, err := EpubFile.AddSection(strings.Join(imgTags, "\n"), chapterTitle, "", "")
	if err != nil {
		return "", err
	}

	return destination, nil
}

// PackToCBZ packs chapter to .cbz format
func PackToCBZ(images []*bytes.Buffer, destination string, context *PackerContext) (string, error) {

	zipFile, err := PackToZip(images, destination, context)
	if err != nil {
		return "", err
	}

	// replace .zip extension with .cbz
	cbzFile := strings.TrimSuffix(zipFile, filepath.Ext(zipFile)) + ".cbz"

	err = RemoveIfExists(cbzFile)
	if err != nil {
		return "", err
	}

	// rename .zip file to .cbz file
	if err := Afero.Rename(zipFile, cbzFile); err != nil {
		return "", err
	}
	return cbzFile, nil
}

func generateComicInfo(context *PackerContext) string {
	return `
<ComicInfo xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <Title>` + context.Chapter.Info + `</Title>
  <Series>` + context.Manga.Info + `</Series>
  <Number>` + strconv.Itoa(context.Chapter.Index) + `</Number>
  <Notes></Notes>
  <Genre>Web Comic</Genre>
  <Web>` + context.Manga.Address + `</Web>
  <Manga>Yes</Manga>
</ComicInfo>
`
}

// PackToZip packs chapter to .zip format
func PackToZip(images []*bytes.Buffer, destination string, context *PackerContext) (string, error) {
	destination += ".zip"
	err := RemoveIfExists(destination)
	if err != nil {
		return "", err
	}

	// Create parent directory since zip have some troubles when it doesn't exist
	if exists, err := Afero.Exists(filepath.Dir(destination)); err != nil {
		return "", err
	} else if !exists {
		if err = Afero.MkdirAll(filepath.Dir(destination), 0777); err != nil {
			return "", err
		}
	}

	// Create zip file
	zipFile, err := Afero.Create(destination)
	if err != nil {
		return "", err
	}
	defer func(zipFile afero.File) {
		_ = zipFile.Close()
	}(zipFile)

	// Create zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		_ = zipWriter.Close()
	}(zipWriter)

	// Add images to zip file
	for i, image := range images {
		if err = addFileToZip(zipWriter, image, PadZeros(i, 4)+".jpg"); err != nil {
			return "", err
		}
	}

	// add ComicInfo.xml to zip file if this is called from CBZ packer
	if context != nil && UserConfig.Formats.Default == CBZ && UserConfig.Formats.Comicinfo {
		// not worth raising an error if this fails
		_ = addFileToZip(zipWriter, bytes.NewBufferString(generateComicInfo(context)), "ComicInfo.xml")
	}

	return destination, nil
}

// addFileToZip adds files to zip writer
func addFileToZip(zipWriter *zip.Writer, file *bytes.Buffer, name string) error {
	// Create a new zip file entry
	header := &zip.FileHeader{
		Name:   name,
		Method: zip.Store,
	}

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

// PackToPlain packs chapter in the plain format
func PackToPlain(images []*bytes.Buffer, destination string, context *PackerContext) (string, error) {
	err := Afero.MkdirAll(destination, 0700)
	if err != nil {
		return "", err
	}

	// add images to folder
	for i, image := range images {
		if err = Afero.WriteFile(
			filepath.Join(destination, PadZeros(i, 4)+".jpg"),
			image.Bytes(),
			0777,
		); err != nil {
			return "", err
		}
	}

	return destination, nil
}

// Packers is a list of packers for available formats
var Packers = map[FormatType]func([]*bytes.Buffer, string, *PackerContext) (string, error){
	PDF:   PackToPDF,
	Plain: PackToPlain,
	Zip:   PackToZip,
	CBZ:   PackToCBZ,
	Epub:  PackToEpub,
}
