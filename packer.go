package main

import (
	"archive/zip"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/afero"
	"io"
	"path/filepath"
	"strconv"
)

func PackToPDF(images []string, destination string) (string, error) {
	destination += ".pdf"

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

func PackToCBZ(images []string, destination string) (string, error) {
	zipFile, err := PackToZip(images, destination, false)
	if err != nil {
		return "", err
	}

	cbzFile := zipFile + ".cbz"

	if err := Afero.Rename(zipFile, cbzFile); err != nil {
		return "", err
	}
	return cbzFile, nil
}

func PackToZip(images []string, destination string, withExtension bool) (string, error) {
	if withExtension {
		destination += ".zip"
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
}
