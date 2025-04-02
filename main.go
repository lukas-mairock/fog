package main

//  Author: Lukas Mairock
//  Date:   2025-01-10

import (
	"flag"
	"fmt"
	"os"

	"luks.cat/src/fetcher"
	"luks.cat/src/renderer"
)

// --------------------------------------------------------------------------------------------------------------------

var (
	version     = "0.1.0-alpha"
	showHelp    = flag.Bool("help", false, "Show help information")
	showVersion = flag.Bool("version", false, "Show version information")
)

// --------------------------------------------------------------------------------------------------------------------

func printHelp() {
	fmt.Println("Usage of", os.Args[0]+":")
	flag.PrintDefaults()
}

// --------------------------------------------------------------------------------------------------------------------

func main() {
	flag.Parse()

	if *showHelp {
		printHelp()
		return
	}

	if *showVersion {
		fmt.Println("Version:", version)
		return
	}

	renderer.Render(*fetcher.Fetch())
}
