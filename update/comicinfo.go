package update

import (
	"encoding/xml"
	"fmt"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"os"
	"path/filepath"
	"strings"
)

func getAnyChapterComicInfo(mangaPath string) (*source.ComicInfo, error) {
	// recursively search for .cbz files
	// find the first one and get the name from it
	var cbzFiles []string
	err := filepath.Walk(mangaPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".cbz") {
			cbzFiles = append(cbzFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(cbzFiles) == 0 {
		return nil, fmt.Errorf("no .cbz files found")
	}

	comicInfo, err := getComicInfoXML(cbzFiles[0])
	if err != nil {
		return nil, err
	}

	return comicInfo, nil
}

func getComicInfoXML(chapter string) (*source.ComicInfo, error) {
	if !strings.HasSuffix(chapter, ".cbz") {
		return nil, fmt.Errorf("chapter must be a .cbz file")
	}

	// open chapter as ReaderAt
	file, err := filesystem.Api().Open(chapter)
	if err != nil {
		return nil, err
	}

	filesystem.SetMemMapFs()
	defer filesystem.SetOsFs()

	// extract ComicInfo.xml
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// No need to delete the file as it will be deleted when the memmap filesystem is reset
	err = util.Unzip(file, stat.Size(), "T")
	if err != nil {
		return nil, err
	}

	// read ComicInfo.xml
	contents, err := filesystem.Api().ReadFile(filepath.Join("T", "ComicInfo.xml"))
	if err != nil {
		return nil, err
	}

	// parse ComicInfo.xml
	var comicInfo source.ComicInfo
	err = xml.Unmarshal(contents, &comicInfo)
	if err != nil {
		return nil, err
	}

	return &comicInfo, nil

}
