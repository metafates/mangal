package update

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

	comicInfo, err := getAnyChapterComicInfo(manga)
	if err != nil {
		return "", err
	}

	nameCache[manga] = comicInfo.Series
	return comicInfo.Series, nil
}
