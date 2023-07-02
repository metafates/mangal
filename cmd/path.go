package cmd

import (
	"fmt"
	"github.com/mangalorg/mangal/path"
)

type pathCmd struct {
	Cache        bool `help:"Path to cache directory"`
	Temp         bool `help:"Path to temporary directory"`
	Downloads    bool `help:"Path to downloads directory"`
	LuaProviders bool `help:"Path to lua providers directory"`
	Header       bool `help:"Print header" negatable:"" default:"true"`
}

func (p *pathCmd) Run() error {
	paths := []struct {
		Name  string
		Func  func() string
		Print bool
	}{
		{"Cache", path.CacheDir, p.Cache},
		{"Temp", path.TempDir, p.Temp},
		{"Downloads", path.DownloadsDir, p.Downloads},
		{"Lua Providers", path.LuaProvidersDir, p.LuaProviders},
	}

	var anyPrinted bool
	for _, t := range paths {
		if t.Print {
			anyPrinted = true
			if p.Header {
				fmt.Println(t.Name)
			}
			fmt.Println(t.Func())
		}
	}

	if !anyPrinted {
		for _, t := range paths {
			if p.Header {
				fmt.Println(t.Name)
			}
			fmt.Println(t.Func())
		}
	}

	return nil
}
