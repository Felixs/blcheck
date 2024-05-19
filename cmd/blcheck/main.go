package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Felixs/blcheck/pkg/url"
)

var (
	flagUrl          string
	processStarttime time.Time
)

// blcheck - A simple tool to check which links on your websites are broken.
func main() {
	// flagUrl := ""
	processStarttime = time.Now()
	checkedArguments(&flagUrl)

	processUrl(flagUrl)
}

// Parses the command line arguments and checks if they all needed arguments are present.
func checkedArguments(flagUrl *string) {
	flag.Parse()
	if flag.NArg() != 1 {
		printUsage("URL is required")
		os.Exit(3)

	}
	*flagUrl = flag.Arg(0)
}

// Prints how to use the tool to stdout with an error message.
func printUsage(errorMsg string) {
	fmt.Printf(`blcheck - A simple tool to check which links on your websites are broken.
	
Usage: blcheck <URL>
	
Error: %s
`, errorMsg)
}

// Checks url for broken links.
func processUrl(inputUrl string) {
	log.Println("Checking URL: ", inputUrl)
	url.InferHttpsPrefix(&inputUrl)
	if !url.IsUrlValid(inputUrl) {
		fmt.Printf("not a valid url %s\n", inputUrl)
		os.Exit(1)
	}

	body, err := url.GetBodyFromUrl(inputUrl)
	if err != nil {
		fmt.Printf("could not get data from url: " + err.Error())
		os.Exit(2)
	}
	httpUrls := url.ExtractHttpUrls(body)
	printUrlOverview(httpUrls)
	fmt.Println("String url scan")
	for i, httpUrl := range httpUrls {
		isAvailbale := url.UrlIsAvailable(httpUrl)
		fmt.Printf("%v %s %d/%d\n", isAvailbale, httpUrl, i+1, len(httpUrls))
	}
	fmt.Printf("Process done in %.4f seconds\n", time.Since(processStarttime).Seconds())
}

// Prints basic infos about the found urls.
func printUrlOverview(httpUrls []string) {
	for _, httpUrl := range httpUrls {
		fmt.Println("Found url " + httpUrl)
	}
	fmt.Printf("Found a total of %d unique urls\n", len(httpUrls))

	fmt.Printf("Process done in %.4f seconds\n", time.Since(processStarttime).Seconds())
}
