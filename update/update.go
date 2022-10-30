package update

import (
	"github.com/metafates/mangal/source"
)

func Chapter(path string, with *source.Chapter) error {
	return nil
}

func Manga(path string, with *source.Manga) error {
	name, err := GetName(path)
	if err != nil {
		return err
	}

	manga := &source.Manga{
		Name: name,
	}

	err = manga.PopulateMetadata(func(string) {})
	if err != nil {
		return err
	}

	return nil
}
