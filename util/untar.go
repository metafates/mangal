package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/metafates/mangal/filesystem"
	"io"
	"os"
	"path/filepath"
)

func UntarGZ(gzipStream io.Reader, path string) error {
	err := filesystem.Api().MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		name := filepath.Join(path, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err = filesystem.Api().MkdirAll(name, os.ModePerm); err != nil {
				return err
			}
		case tar.TypeReg:
			err = filesystem.Api().MkdirAll(filepath.Dir(name), os.ModePerm)
			if err != nil {
				return err
			}

			outFile, err := filesystem.Api().OpenFile(name, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				return err
			}
			if _, err = io.Copy(outFile, tarReader); err != nil {
				return err
			}
			err = outFile.Close()
		default:
			err = fmt.Errorf(
				"uknown type: %s in %s",
				string(header.Typeflag),
				header.Name)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
