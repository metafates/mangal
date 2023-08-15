package config

//go:generate enumer -type=Mode -trimprefix=Mode -json -yaml -text
type Mode uint8

const (
	ModeNone Mode = iota + 1
	ModeTUI
	ModeWeb
	ModeScript
)
