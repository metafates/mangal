package mini

import (
	"errors"
	"github.com/metafates/mangal/util"
	"github.com/samber/lo"
)

var (
	pageSize = 15
	trimAt   = 30
)

type Options struct {
	Download bool
	Continue bool
}

func init() {
	if w, _, err := util.TerminalSize(); err == nil {
		trimAt = lo.Max([]int{trimAt, w - 10})
	}
}

func Run(options *Options) error {
	if options.Continue && options.Download {
		return errors.New("cannot download and continue")
	}

	if options.Continue {
		return continueReading()
	} else if options.Download {
		return download()
	} else {
		return read()
	}
}
