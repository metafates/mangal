package main

import (
	"fmt"
	"github.com/metafates/mangai/cmd"
	"os"
)

// Set on compile time
var (
	version string
	build   string
)

func main() {
	err := cmd.Execute(version, build)
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
