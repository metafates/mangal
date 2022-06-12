package main

import "testing"

func TestNewBubble(t *testing.T) {
	initConfig("")

	bubble := NewBubble(searchState)

	if bubble.state != searchState {
		t.Error("Invalid state")
	}

	if bubble.mangaList.FilteringEnabled() {
		t.Error("Filtered should be disabled for manga list")
	}

	if bubble.chaptersList.FilteringEnabled() {
		t.Error("Filtered should be disabled for chapters list")
	}

	if bubble.mangaChan == nil {
		t.Error("Manga channel should not be nil")
	}

	if bubble.chaptersChan == nil {
		t.Error("Chapters channel should not be nil")
	}

	if bubble.chaptersProgressChan == nil {
		t.Error("Chapters progress channel should not be nil")
	}

	if bubble.chapterPagesProgressChan == nil {
		t.Error("Chapter pages progress channel should not be nil")
	}

	if !bubble.input.Focused() {
		t.Error("Input should be focused")
	}
}
