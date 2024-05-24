/*
	blcheck - A simple tool to check which links on your websites are broken.

Usage: blcheck <URL>

-max-parallel-requests int

	Setting a maximum how many requests get executed in parallel (default 20)

-max-response-timeout int

	Maximum timeout wait on requests in seconds (default 5)

-version

	Displays version of blcheck

-json

	Export output as json format
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

const version = "0.0.2"

var (
	flagVersion             bool
	flagUrl                 string
	flagMaxParallelRequests int
	flagMaxTimeoutInSeconds int
	flagOutputJson          bool
	errorMessage            string
)

// Blcheck entry point.
func main() {
	checkedArguments(&flagUrl)
	processUrl(flagUrl)
}

// Parses the command line arguments and checks if they all needed arguments are present.
func checkedArguments(flagUrl *string) {
	// TODO: if using 2 different flag names for a single values flag.Usage needs to be overwritten
	// Version
	flag.BoolVar(&flagVersion, "version", false, "Displays version of blcheck")
	flag.BoolVar(&flagVersion, "v", false, "Displays version of blcheck")
	// Ratelimit / parallel requests
	flag.IntVar(&flagMaxParallelRequests, "max-parallel-requests", report.MaxNumParallelQueries, "Maximum number of parallel requests executed")
	flag.IntVar(&flagMaxParallelRequests, "mpr", report.MaxNumParallelQueries, "Maximum number of parallel requests executed")
	// Timeout
	flag.IntVar(&flagMaxTimeoutInSeconds, "max-response-timeout", int(url.DefaultHttpGetTimeout.Seconds()), "Maximum timeout wait on requests in seconds")
	flag.IntVar(&flagMaxTimeoutInSeconds, "mrt", int(url.DefaultHttpGetTimeout.Seconds()), "Maximum timeout wait on requests in seconds")
	// Output as json flag
	flag.BoolVar(&flagOutputJson, "json", false, "Export output as json format")
	flag.BoolVar(&flagOutputJson, "j", false, "Export output as json format")

	// setting own print function, to handle positonal arguments
	flag.Usage = printUsage
	flag.Parse()

	if flagVersion {
		fmt.Println("blcheck " + version + "\n2024 - Felix Sponholz")
		os.Exit(0)
	}
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

	fmt.Printf(`blcheck (%s)- A simple tool to check which links on your websites are broken.
	
Usage: blcheck <URL>
`, version)
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
	parsing_duration := time.Since(processStarttime).String()

	url.SetHttpGetTimeoutSeconds(time.Duration(flagMaxTimeoutInSeconds) * time.Second)
	urlReports := report.CustomizableCreateUrlReport(httpUrls, int(flagMaxParallelRequests))
	urlReports.AddMetaData("initial_parsing_duration", parsing_duration)

	switch {
	case flagOutputJson:
		fmt.Println(urlReports.Json())
	default:
		fmt.Println(urlReports.FullString())
	}

}
