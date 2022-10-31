package update

import (
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"path/filepath"
)

func getSeriesJSON(manga string) (*source.SeriesJSON, error) {
	// check if series.json exists at manga dir
	serisJSONPath := filepath.Join(manga, "series.json")
	exists, err := filesystem.Api().Exists(serisJSONPath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("series.json must be present")
	}

	contents, err := filesystem.Api().ReadFile(serisJSONPath)
	if err != nil {
		return nil, err
	}

	var seriesJSON source.SeriesJSON

	err = json.Unmarshal(contents, &seriesJSON)
	return &seriesJSON, err
}
