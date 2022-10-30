package update

import (
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/anilist"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"path/filepath"
)

func Chapter(path string, with *source.Chapter) error {
	return nil
}

func Manga(path string, with *source.Manga) error {
	return nil
}

func getName(manga string) (string, error) {
	// check if series.json exists at manga dir
	serisJSONPath := filepath.Join(manga, "series.json")
	exists, err := filesystem.Api().Exists(serisJSONPath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return "", fmt.Errorf("series.json must be present")
	}

	contents, err := filesystem.Api().ReadFile(serisJSONPath)
	if err != nil {
		return "", err
	}

	var seriesJSON source.SeriesJSON

	err = json.Unmarshal(contents, &seriesJSON)
	if err != nil {
		return "", err
	}

	return seriesJSON.Metadata.Name, nil
}

func asAnilist(manga string) (*anilist.Manga, error) {
	name, err := getName(manga)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
