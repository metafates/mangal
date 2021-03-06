package main

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestInlineMode(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping inline mode")
	}

	// test options
	options := InlineOptions{
		config:     "",
		mangaIdx:   -1,
		chapterIdx: -1,
		asJson:     false,
		format:     PDF,
		showUrls:   false,
		asTemp:     true,
		doRead:     false,
		doOpen:     false,
	}

	out, err := InlineMode(TestQuery, options)
	if err != nil {
		t.Error(err)
	}

	if out == "" {
		t.Error("output is empty")
	}

	// changet test options
	options.asJson = true

	out, err = InlineMode(TestQuery, options)

	if err != nil {
		t.Error(err)
	}

	if out == "" {
		t.Error("output is empty")
	}

	// parse json
	var result []*URL
	err = json.Unmarshal([]byte(out), &result)
	if err != nil {
		t.Error(err)
	}

	if len(result) == 0 {
		t.Error("result is empty")
	}

	// select invalid manga
	options.mangaIdx = 0

	out, err = InlineMode(TestQuery, options)
	if err == nil {
		t.Error("expected error")
	}

	// select invalid chapter
	options.chapterIdx = 0

	out, err = InlineMode(TestQuery, options)
	if err == nil {
		t.Error("expected error")
	}

	// select valid manga
	options.mangaIdx = 1
	options.chapterIdx = -1

	out, err = InlineMode(TestQuery, options)
	if err != nil {
		t.Error(err)
	}

	if out == "" {
		t.Error("output is empty")
	}

	// select valid chapter
	options.chapterIdx = 1

	out, err = InlineMode(TestQuery, options)
	if err != nil {
		t.Error(err)
	}

	if out == "" {
		t.Error("output is empty")
	}

	// download as temp
	options.asTemp = true

	out, err = InlineMode(TestQuery, options)
	if err != nil {
		t.Error(err)
	}

	if out == "" {
		t.Error("output is empty")
	}

	// check if file at out path exists
	if _, err = Afero.Exists(out); err != nil {
		t.Error(err)
	}

	// check if file at out path is not empty
	if _, err = Afero.Stat(out); err != nil {
		t.Error(err)
	}

	// check file extension
	if filepath.Ext(out) != ".pdf" {
		t.Error(err)
	}
}
