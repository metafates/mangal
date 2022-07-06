package tui

import (
	"github.com/metafates/mangal/config"
	"testing"
)

func TestNewBubble(t *testing.T) {
	config.Initialize("", false)

	bubble := NewBubble(SearchState)

	if bubble.state != SearchState {
		t.Error("Invalid state")
	}

	if bubble.mangaList.FilteringEnabled() {
		t.Error("Filtered should be disabled for manga list")
	}

	if bubble.chaptersList.FilteringEnabled() {
		t.Error("Filtering should be disabled for chapters list")
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

	if bubble.input.Value() != "" {
		t.Error("Input should be empty")
	}
}
