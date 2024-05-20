/*
	blcheck - A simple tool to check which links on your websites are broken.

Usage: blcheck <URL>

	-max-parallel-requests int
	      Setting a maximum how many requests get executed in parallel (default 20)
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Felixs/blcheck/pkg/report"
	"github.com/Felixs/blcheck/pkg/url"
)

var (
	flagUrl                 string
	flagMaxParallelRequests uint
	errorMessage            string
)

// Blcheck entry point.
func main() {
	checkedArguments(&flagUrl)
	processUrl(flagUrl)
}

// Parses the command line arguments and checks if they all needed arguments are present.
func checkedArguments(flagUrl *string) {
	flag.UintVar(&flagMaxParallelRequests, "max-parallel-requests", report.MaxNumParallelQueries, "Setting a maximum how many requests get executed in parallel")
	// setting own print function, to handle positonal arguments
	flag.Usage = printUsage
	flag.Parse()
	if flag.NArg() != 1 {
		errorMessage = "URL is required"
		printUsage()
		os.Exit(3)

	}
	*flagUrl = flag.Arg(0)
}

// Prints how to use the tool to stdout, with an error message if present.
func printUsage() {
	if errorMessage != "" {
		fmt.Println(errorMessage)
	}

	fmt.Printf(`blcheck - A simple tool to check which links on your websites are broken.
	
Usage: blcheck <URL>
`)
	flag.PrintDefaults()
}

// Checks url for broken links.
func processUrl(inputUrl string) {
	log.Println("Checking URL: ", inputUrl)
	processStarttime := time.Now()
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
	fmt.Printf("parsing done in %.4f seconds\n", time.Since(processStarttime).Seconds())
	fmt.Printf("extracted %d unique urls, starting url report scan\n", len(httpUrls))

	urlReports := report.CustomizableCreateUrlReport(httpUrls, int(flagMaxParallelRequests))
	fmt.Println(urlReports.FullString())

}
