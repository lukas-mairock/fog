package renderer

import (
	"fmt"
	"strings"

	"luks.cat/src/fetcher"
)

// --------------------------------------------------------------------------------------------------------------------

func getPadding(text string, padstring string) string {
	var targetLength int = 48
	var padding string
	if len(text) < targetLength {
		padding = strings.Repeat(padstring, targetLength-len(text))
	}
	return padding
}

// --------------------------------------------------------------------------------------------------------------------

func printHeader(text string) {
	color := "\x1b[0;34m"
	padding := getPadding(fmt.Sprintf("╭───── %s ",text), "─")               
	fmt.Printf("%s╭────── \x1b[0;1m%s%s %s\x1b[0m\n", color, text, color, padding)
}

// --------------------------------------------------------------------------------------------------------------------

func printSeperator(text string) {
	color := "\x1b[0;34m"
	padding := getPadding(fmt.Sprintf("├───── %s ",text), "─")               
	fmt.Printf("%s├────── \x1b[0;1m%s%s %s\x1b[0m\n", color, text, color, padding)
}

// --------------------------------------------------------------------------------------------------------------------

func printFooter() {
	color := "\x1b[0;34m"
	padding := getPadding(fmt.Sprintf("╰─────"), "─")     
	fmt.Printf("%s╰──────%s\x1b[0m\n", color, padding)
}

// --------------------------------------------------------------------------------------------------------------------

func printItem(key string, value string) {
	fstring := fmt.Sprintf("\x1b[34m│\x1b[0;1;32m %s:\x1b[0m\t%s", key, value)
	fmt.Println(fstring)
}

// --------------------------------------------------------------------------------------------------------------------

func Render(info fetcher.Information) {
	printHeader("System")
	printItem("Distro", info.Sys.Distribution)
	printItem("Hostname", info.Sys.Hostname)
	printItem("Address", info.Sys.Address)
	printItem("Uptime", info.Sys.Uptime)
	printItem("Kernel", info.Sys.Kernel)

	if info.Cron != nil {
		printSeperator("Crons")
		for _, job := range info.Cron {
			printItem(job.Description, job.Schedule)
		}
	}

	if info.Task != nil {
		printSeperator("Tasks")
		for index, task := range info.Task.Lines {
			printItem(fmt.Sprintf("%d", index+1), fmt.Sprintf("\t%s", task))
		}
	}

	printFooter()
}
