package info

import (
	"io"

	"github.com/mangalorg/libmangal"
	"github.com/pelletier/go-toml"
)

//go:generate enumer -type=Type -trimprefix=Type -json -text
type Type uint8

const (
	TypeBundle Type = iota + 1
	TypeLua
)

const Filename = "mangal.toml"

// Info contains libmangal info about provider with mangal specific type field
type Info struct {
	Info libmangal.ProviderInfo
	Type Type
}

// Parse parses info from reader
func Parse(r io.Reader) (info Info, err error) {
	decoder := toml.NewDecoder(r)
	decoder.Strict(true)

	err = decoder.Decode(&info)
	return
}
