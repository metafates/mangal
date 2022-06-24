package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// MockFiles creates a temporary directory and fills it with files.
// It is caller responsibility to remove created files.
func MockFiles(t *testing.T) []string {
	t.Helper()

	// options to get chapter images
	options := InlineOptions{
		config:     "",
		mangaIdx:   1,
		chapterIdx: 1,
		asJson:     false,
		format:     Plain,
		showUrls:   false,
		asTemp:     true,
		doRead:     false,
		doOpen:     false,
	}

	// get chapter images
	out, err := InlineMode(TestQuery, options)
	if err != nil {
		t.Fatal(err)
	}

	// get files in downloaded folder from output
	images, err := Afero.ReadDir(out)

	if err != nil {
		t.Fatal(err)
	}

	// get images from the chapter
	if len(images) == 0 {
		t.Fatal("no images created")
	}

	return Map(images, func(i os.FileInfo) string {
		return filepath.Join(out, i.Name())
	})
}

// TempFile creates a temporary file with given suffix.
// It is caller responsibility to remove created file.
func TempFile(t *testing.T, extension string) string {
	t.Helper()

	// create temp file
	f, err := Afero.TempFile(os.TempDir(), TempPrefix+"*"+extension)
	if err != nil {
		t.Fatal(err)
	}

	// close file
	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	return filepath.Join(os.TempDir(), f.Name())
}

func TestPackToPDF(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping this PackToPDF is too expensive")
	}

	// create mock files
	files := MockFiles(t)

	if len(files) == 0 {
		t.Fatal("no files created")
	}

	// pack to pdf
	out, err := PackToPDF(files, TempFile(t, ".pdf"))

	if err != nil {
		t.Fatal(err)
	}

	if out == "" {
		t.Fatal("output is empty")
	}

	// check if output is a pdf
	if !strings.HasSuffix(out, ".pdf") {
		t.Error("output is not a pdf")
	}

	// check if pdf is not empty
	if stat, err := Afero.Stat(out); err != nil {
		t.Fatal(err)
	} else if stat.Size() == 0 {
		t.Error("pdf is empty")
	}

	// remove mock files
	RemoveTemp()
}

func TestPackToCBZ(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping this PackToCBZ is too expensive")
	}

	// create mock files
	files := MockFiles(t)

	if len(files) == 0 {
		t.Fatal("no files created")
	}

	// pack to cbz
	out, err := PackToCBZ(files, TempFile(t, ".cbz"))

	if err != nil {
		t.Fatal(err)
	}

	if out == "" {
		t.Fatal("output is empty")
	}

	// check if output is a cbz
	if !strings.HasSuffix(out, ".cbz") {
		t.Error("output is not a cbz")
	}

	// check if cbz is not empty
	if stat, err := Afero.Stat(out); err != nil {
		t.Fatal(err)
	} else if stat.Size() == 0 {
		t.Error("cbz is empty")
	}

	// remove mock files
	err = Afero.RemoveAll(filepath.Dir(files[0]))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPackToEpub(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping this PackToEPUB is too expensive")
	}

	// create mock files
	files := MockFiles(t)

	if len(files) == 0 {
		t.Fatal("no files created")
	}

	// pack to epub
	out, err := PackToEpub(files, TempFile(t, ".epub"))

	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}

	if out == "" {
		t.Fatal("output is empty")
	}

	err = EpubFile.Write(out)

	if err != nil {
		t.Fatal(err)
	}

	// check if output is a epub
	if !strings.HasSuffix(out, ".epub") {
		t.Error("output is not a epub")
	}

	// check if epub is not empty
	if stat, err := Afero.Stat(out); err != nil {
		t.Fatal(err)
	} else if stat.Size() == 0 {
		t.Error("epub is empty")
	}

	// remove mock files
	err = Afero.RemoveAll(filepath.Dir(files[0]))
	if err != nil {
		t.Fatal(err)
	}
}
