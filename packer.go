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

// PackToPDF packs chapter to .pdf format
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

	if err = pdfcpu.ImportImagesFile(images, destination, nil, nil); err != nil {
		return "", err
	}

	return destination, nil
}

// EpubFile global epub file that is used to save multiple chapters in a single file
var EpubFile *epub.Epub

// PackToEpub adds chapter to epub file
func PackToEpub(images []string, destination string) (string, error) {
	coverSet := true
	mangaTitle := filepath.Base(filepath.Dir(destination))
	chapterTitle := filepath.Base(destination)

	destination = filepath.Join(filepath.Dir(destination), mangaTitle+".epub")

	// Initialize epub file
	if EpubFile == nil {
		coverSet = false
		EpubFile = epub.NewEpub(mangaTitle)
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
		epubImage, err := EpubFile.AddImage(image, strconv.Itoa(rand.Intn(100000))+filepath.Base(image))
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
func PackToCBZ(images []string, destination string) (string, error) {

	zipFile, err := PackToZip(images, destination)
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

// PackToZip packs chapter to .zip format
func PackToZip(images []string, destination string) (string, error) {
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
		if err = addFileToZip(zipWriter, image, i); err != nil {
			return "", err
		}
	}

	return destination, nil
}

// addFileToZip adds files to zip writer
func addFileToZip(zipWriter *zip.Writer, filename string, index int) error {
	// Open file
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

	// Create a new zip file entry
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Set the desired header name
	header.Name = PadZeros(index, 4) + filepath.Ext(header.Name)

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

	// add images to folder
	for i, image := range images {
		imageContents, err := Afero.ReadFile(image)
		if err != nil {
			return "", err
		}

		if err = Afero.WriteFile(
			filepath.Join(destination, PadZeros(i, 4)+filepath.Ext(image)),
			imageContents,
			0700,
		); err != nil {
			return "", err
		}
	}

	return destination, nil
}

// Packers is a list of packers for available formats
var Packers = map[FormatType]func([]string, string) (string, error){
	PDF:   PackToPDF,
	Plain: PackToPlain,
	Zip:   PackToZip,
	CBZ:   PackToCBZ,
	Epub:  PackToEpub,
}
