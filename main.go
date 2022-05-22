package main

import (
	"fmt"
	"github.com/metafates/mangai/cmd"
	"math/rand"
	"os"
	"time"
)

// Set on compile time
var (
	version string
	build   string
)

func main() {
	rand.Seed(time.Now().Unix())
	err := cmd.Execute(version, build)
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
