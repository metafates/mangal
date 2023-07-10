package cmd

import (
	"fmt"
	"github.com/mangalorg/mangal/meta"
)

type versionCmd struct {
	Short bool `help:"just show the version number"`
}

func (v *versionCmd) Run() error {
	if v.Short {
		fmt.Println(meta.Version)
		return nil
	}

	fmt.Println(meta.PrettyVersion())
	return nil
}
