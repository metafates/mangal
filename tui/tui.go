package tui

import "errors"

type Options struct {
	Continue bool
}

func Run(options *Options) error {
	return errors.New("tui mode is not implemented yet. use mini")
}
