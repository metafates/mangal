package util

import (
	"archive/zip"
	"fmt"
	"github.com/metafates/mangal/filesystem"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(zipStream io.ReaderAt, size int64, dest string) error {
	r, err := zip.NewReader(zipStream, size)
	if err != nil {
		return err
	}

	err = filesystem.Api().MkdirAll(dest, os.ModePerm)
	if err != nil {
		return err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		defer Ignore(rc.Close)

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			err = filesystem.Api().MkdirAll(path, f.Mode())
			if err != nil {
				return err
			}
		} else {
			err = filesystem.Api().MkdirAll(filepath.Dir(path), f.Mode())
			if err != nil {
				return err
			}

			f, err := filesystem.Api().OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			defer Ignore(f.Close)

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)

		if err != nil {
			return err
		}
	}

	return nil
}
