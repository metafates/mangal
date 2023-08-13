package config

//go:generate enumer -type=Mode -trimprefix=Mode -json -yaml -text
type Mode uint8

const (
	ModeTUI Mode = iota + 1
	ModeWeb
	ModeScript
)
