package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Felixs/blcheck/pkg/url"
)

// blcheck - A simple tool to check which links on your websites are broken
func main() {
	flagUrl := ""
	checkedArguments(&flagUrl)

	processUrl(flagUrl)
}

// parses the command line arguments and checks if they all needed arguments are present
func checkedArguments(flagUrl *string) {
	flag.Parse()
	if flag.NArg() != 1 {
		printUsage("URL is required")
		os.Exit(3)

	}
	*flagUrl = flag.Arg(0)
}

// prints how to use the tool to stdout with an error message
func printUsage(errorMsg string) {
	fmt.Printf(`blcheck - A simple tool to check which links on your websites are broken.
	
Usage: blcheck <URL>
	
Error: %s
`, errorMsg)
}

// checks url for broken links
func processUrl(input_url string) {
	log.Println("Checking URL: ", input_url)
	if !url.IsUrlValid(input_url) {
		fmt.Printf("not a valid url %s\n", input_url)
		os.Exit(1)
	}
	body, err := url.GetBodyFromUrl(input_url)
	if err != nil {
		fmt.Printf("could not get data from url: " + err.Error())
		os.Exit(2)
	}
	fmt.Println(body)
}
