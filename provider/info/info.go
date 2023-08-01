package info

import "github.com/mangalorg/libmangal"

//go:generate enumer -type=Type -trimprefix=Type -json -text
type Type uint8

const (
	TypeLua Type = iota + 1
	TypeBundle
)

const Filename = "mangal.toml"

type Info struct {
	Info libmangal.ProviderInfo
	Type Type
}
