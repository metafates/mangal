package config

import (
	"github.com/metafates/mangal/common"
	"testing"
)

func TestParseConfig(t *testing.T) {
	config, err := ParseConfig([]byte(common.DefaultConfigString))

	if err != nil {
		t.Fatal(err)
	}

	if config.UI.Fullscreen == false {
		t.Error("Fullscreen is false")
	}

	if config.Downloader.Path != "." {
		t.Error("Downloader.Path is not .")
	}

	if config.Downloader.CacheImages == true {
		t.Error("Downloader.CacheImages is true")
	}

	if config.Anilist.Enabled == true {
		t.Error("Anilist.Enabled is true")
	}

	if config.Anilist.MarkDownloaded == true {
		t.Error("Anilist.MarkDownloaded is true")
	}

	if config.UseCustomReader == true {
		t.Error("UseCustomReader is true")
	}

	if config.Formats.Default != PDF {
		t.Error("Formats.Default is not PDF")
	}
}

func TestGetConfig(t *testing.T) {
	config := GetConfig("")

	if err := ValidateConfig(config); err != nil {
		t.Error(err)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Fatal("Error while parsing default config file")
	}

	if err := ValidateConfig(config); err != nil {
		t.Error(err)
	}
}
