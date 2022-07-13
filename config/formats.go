package config

import "github.com/metafates/mangal/common"

type FormatsConfig struct {
	Default   common.FormatType `toml:"default"`
	Comicinfo bool
}
