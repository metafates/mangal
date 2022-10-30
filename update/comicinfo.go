package update

import (
	"encoding/xml"
	"fmt"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"path/filepath"
	"strings"
)

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
