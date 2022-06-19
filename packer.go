package main

import (
	"archive/zip"
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

// PackToPDF packs chapter to .pdf
func PackToPDF(images []string, destination string) (string, error) {
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

	if err := pdfcpu.ImportImagesFile(images, destination, nil, nil); err != nil {
		return "", err
	}

	return destination, nil
}

// EpubFile global epub file that is used to save multiple chapters in a single file
var EpubFile *epub.Epub

// PackToEpub adds chapter to the epub file
func PackToEpub(images []string, destination string) (string, error) {
	mangaTitle := filepath.Base(filepath.Dir(destination))
	chapterTitle := filepath.Base(destination)

	destination = filepath.Join(filepath.Dir(destination), mangaTitle+".epub")

	// Initialize epub file
	if EpubFile == nil {
		EpubFile = epub.NewEpub(mangaTitle)
		EpubFile.SetPpd("rtl")
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

	for i, image := range images {
		epubImage, err := EpubFile.AddImage(image, strconv.Itoa(rand.Intn(100000))+filepath.Base(image))
		if err != nil {
			return "", err
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
func PackToCBZ(images []string, destination string) (string, error) {

	zipFile, err := PackToZip(images, destination, false)
	if err != nil {
		return "", err
	}

	cbzFile := zipFile + ".cbz"

	err = RemoveIfExists(cbzFile)
	if err != nil {
		return "", err
	}

	if err := Afero.Rename(zipFile, cbzFile); err != nil {
		return "", err
	}
	return cbzFile, nil
}

// PackToZip packs chapter to .zip format
func PackToZip(images []string, destination string, withExtension bool) (string, error) {
	if withExtension {
		destination += ".zip"
		err := RemoveIfExists(destination)
		if err != nil {
			return "", err
		}
	}

	if exists, err := Afero.Exists(filepath.Dir(destination)); err != nil {
		return "", err
	} else if !exists {
		if err := Afero.MkdirAll(filepath.Dir(destination), 0777); err != nil {
			return "", err
		}
	}

	zipFile, err := Afero.Create(destination)
	if err != nil {
		return "", err
	}
	defer func(zipFile afero.File) {
		_ = zipFile.Close()
	}(zipFile)

	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		_ = zipWriter.Close()
	}(zipWriter)

	for i, image := range images {
		if err = addFileToZip(zipWriter, image, i); err != nil {
			return "", err
		}
	}

	return destination, nil
}

// addFileToZip adds files to zip writer
func addFileToZip(zipWriter *zip.Writer, filename string, index int) error {
	file, err := Afero.Open(filename)
	if err != nil {
		return err
	}

	defer func(file afero.File) {
		_ = file.Close()
	}(file)

	// Get the file information
	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = strconv.Itoa(index) + filepath.Ext(header.Name)

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

// PackToPlain packs chapter in the plain format
func PackToPlain(images []string, destination string) (string, error) {
	err := Afero.MkdirAll(destination, 0700)
	if err != nil {
		return "", err
	}

	for i, image := range images {
		imageContents, err := Afero.ReadFile(image)
		if err != nil {
			return "", err
		}

		if err = Afero.WriteFile(
			filepath.Join(destination, strconv.Itoa(i)+filepath.Ext(image)),
			imageContents,
			0700,
		); err != nil {
			return "", err
		}
	}

	return destination, nil
}

var Packers = map[FormatType]func([]string, string) (string, error){
	PDF:   PackToPDF,
	Plain: PackToPlain,
	Zip:   func(f []string, d string) (string, error) { return PackToZip(f, d, true) },
	CBZ:   PackToCBZ,
	Epub:  PackToEpub,
}
