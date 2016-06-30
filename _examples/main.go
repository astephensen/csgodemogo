package main

import (
	"fmt"
	"os"

	"github.com/astephensen/csgodemogo"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [demo.dem]\n", os.Args[0])
		os.Exit(2)
	}

	demofile := csgodemogo.Open(os.Args[1])
	demofile.Header.PrintInfo()
	for 1 == 1 {
		demofile.GetFrame()
	}
}
