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

	demofile.GameEventEmitter = func(gameEvent interface{}) {
		switch gameEvent.(type) {
		case csgodemogo.GameEventRoundStart:
			fmt.Println("Round Started", gameEvent)
		case csgodemogo.GameEventRoundEnd:
			fmt.Println("Round Ended")
		}
	}

	demofile.Header.PrintInfo()
	for demofile.Finished == false {
		demofile.GetFrame()
	}

	fmt.Println("Finished parsing demo file!")
}
