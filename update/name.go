package update

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var nameCache = make(map[string]string)

func GetName(manga string) (string, error) {
	if name, ok := nameCache[manga]; ok {
		return name, nil
	}

	seriesJSON, err := getSeriesJSON(manga)
	if err == nil {
		nameCache[manga] = seriesJSON.Metadata.Name
		return seriesJSON.Metadata.Name, nil
	}

	// recursively search for .cbz files
	// find the first one and get the name from it
	var cbzFiles []string
	err = filepath.Walk(manga, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".cbz") {
			cbzFiles = append(cbzFiles, path)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if len(cbzFiles) == 0 {
		return "", fmt.Errorf("no .cbz files found")
	}

	comicInfo, err := getComicInfoXML(cbzFiles[0])
	if err != nil {
		return "", err
	}

	nameCache[manga] = comicInfo.Series
	return comicInfo.Series, nil
}
