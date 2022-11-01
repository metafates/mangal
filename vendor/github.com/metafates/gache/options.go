package gache

import (
	"time"
)

// Options is a struct that contains options for the cache.
type Options struct {
	// Path to the file to store the cache.
	// If not specified (empty string), the cache will be stored in memory.
	Path string

	// Lifetime is the time duration after which the cache expires.
	// Values below zero are treated as never expiring.
	// Defaults to -1.
	Lifetime time.Duration

	// FileSystem is a filesystem that is used to store the cache file.
	FileSystem FileSystem

	// Encoder is the encoder to use for the cache.
	Encoder Encoder

	// Decoder is the decoder to use for the cache.
	Decoder Decoder

	// ExpirationHook is a function that is called when the cache expires.
	ExpirationHook func()
}
