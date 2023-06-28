package cmd

import (
	"fmt"

	"github.com/mangalorg/mangal/meta"
)

type versionCmd struct {}

func (v *versionCmd) Run() error {
    fmt.Printf("%s %s\n", meta.AppName, meta.Version)
    return nil
}
