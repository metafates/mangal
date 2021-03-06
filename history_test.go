package main

import "testing"

func RemoveHistoryFile(t *testing.T) {
	t.Helper()

	path, err := HistoryFile()

	if err != nil {
		t.Fatal(err)
	}

	err = RemoveIfExists(path)

	if err != nil {
		t.Fatal(err)
	}
}

func SampleChapter(t *testing.T) *URL {
	t.Helper()

	manga := &URL{
		Address: "https://anilist.co/manga/1",
		Info:    "test manga",
	}

	chapter := &URL{
		Address:  "https://anilist.co/chapter/1",
		Index:    1,
		Info:     "test chapter",
		Relation: manga,
		Scraper:  UserConfig.Scrapers[0],
	}

	return chapter
}

func TestReadHistory(t *testing.T) {
	RemoveHistoryFile(t)

	history, err := ReadHistory()

	if err != nil {
		t.Fatal(err)
	}

	if len(history) != 0 {
		t.Error("history is not empty")
	}

	RemoveHistoryFile(t)
}

func TestWriteHistory(t *testing.T) {
	RemoveHistoryFile(t)

	chapter := SampleChapter(t)

	err := WriteHistory(chapter)

	if err != nil {
		t.Fatal(err)
	}

	history, err := ReadHistory()

	if err != nil {
		t.Fatal(err)
	}

	if len(history) != 1 {
		t.Error("history must be of length 1")
	}

	manga := chapter.Relation

	if history[manga.Address] == nil {
		t.Error("history entry must not be nil")
	}

	if history[manga.Address].Manga.Address != manga.Address {
		t.Error("manga address must be equal")
	}

	if history[manga.Address].Manga.Info != manga.Info {
		t.Error("manga info must be equal")
	}

	if history[manga.Address].Chapter.Index != chapter.Index {
		t.Error("chapter index must be equal")
	}

	RemoveHistoryFile(t)
}
