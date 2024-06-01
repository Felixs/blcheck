/*
	blcheck - A simple tool to check which links on your websites are broken.

Usage: blcheck <URL>

-max-parallel-requests int

	Setting a maximum how many requests get executed in parallel (default 5)

-max-response-timeout int

	Maximum timeout wait on requests in seconds (default 5)

-version

	Displays version of blcheck

-json

	Export output as json format

-csv

	Export output as csv format (default if no other format given)

-include

	Parsed urls need to contain this string to get checked

-exclude

	Parsed urls need to not contain this string to get checked

-dry

	Only gets urls from initial webpage and does not check the status of other urls

-out

	Writes output to given location. If directory is given, writes to blcheck.log in directory.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Felixs/blcheck/pkg/url"
)

const (
	version   = "0.0.2"
	goodbyMsg = "Thanks for using blcheck. Feel free to check out the repo at https://github.com/Felixs/blcheck"
)

var (
	// Tool comandline flags
	flagOutputJson    bool
	flagOutputCsv     bool
	flagOutputInFile  string
	flagVersion       bool
	flagExecuteDryRun bool

	flagUrl               string
	flagUrlInclude        string
	flagUrlExclude        string
	flagUrlShowReachables bool

	flagMaxParallelRequests int
	flagMaxTimeoutInSeconds int

	// Error message on flag errors/missmatch
	errorMessage string
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
	flag.IntVar(&flagMaxParallelRequests, "max-parallel-requests", url.MaxNumParallelQueries, "Maximum number of parallel requests executed")
	flag.IntVar(&flagMaxParallelRequests, "mpr", url.MaxNumParallelQueries, "Maximum number of parallel requests executed")
	// Timeout
	flag.IntVar(&flagMaxTimeoutInSeconds, "max-response-timeout", int(url.DefaultHttpGetTimeout.Seconds()), "Maximum timeout wait on requests in seconds")
	flag.IntVar(&flagMaxTimeoutInSeconds, "mrt", int(url.DefaultHttpGetTimeout.Seconds()), "Maximum timeout wait on requests in seconds")
	// Output as json flag
	flag.BoolVar(&flagOutputJson, "json", false, "Export output as json format")
	flag.BoolVar(&flagOutputJson, "j", false, "Export output as json format")
	// Output as csv flag
	flag.BoolVar(&flagOutputCsv, "csv", true, "Export output as csv format (default if no other format given)")
	flag.BoolVar(&flagOutputCsv, "c", true, "Export output as csv format (default if no other format given)")
	// Include flag for which string needs to be present in url to check
	flag.StringVar(&flagUrlInclude, "include", "", "Parsed urls need to contain this string to get checked")
	flag.StringVar(&flagUrlInclude, "in", "", "Parsed urls need to contain this string to get checked")
	// Exclude flag for which string can not be present in url to check
	flag.StringVar(&flagUrlExclude, "exclude", "", "Parsed urls need to not contain this string to get checked")
	flag.StringVar(&flagUrlExclude, "ex", "", "Parsed urls need to not contain this string to get checked")
	// Flag if tool should run in dry mode, only getting links from initial webpage
	flag.BoolVar(&flagExecuteDryRun, "dry", false, "Only gets urls from initial webpage and does not check the status of other urls")
	flag.BoolVar(&flagExecuteDryRun, "d", false, "Only gets urls from initial webpage and does not check the status of other urls")
	// Flag if output should be writen into file, gives path an name of file
	flag.StringVar(&flagOutputInFile, "o", "", "Writes output to given location. If directory is given, writes to blcheck.log in directory.")
	flag.StringVar(&flagOutputInFile, "out", "", "Writes output to given location. If directory is given, writes to blcheck.log in directory.")
	// Flag if reachable urls should be included into the output
	flag.BoolVar(&flagUrlShowReachables, "sr", false, "Includes reachable urls in report")
	flag.BoolVar(&flagUrlShowReachables, "show-reachable", false, "Includes reachable urls in report")

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
	url.InferHttpsPrefix(&inputUrl)
	fmt.Println("Checking URL: ", inputUrl)
	processStarttime := time.Now()

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
	numberUniqueUrls := len(httpUrls)

	// check for exclusion
	if flagUrlExclude != "" {
		httpUrls = url.FilterByExclude(httpUrls, flagUrlExclude)
	}
	// check for inclusion
	if flagUrlInclude != "" {
		httpUrls = url.FilterByInclude(httpUrls, flagUrlInclude)
	}

	parsing_duration := time.Since(processStarttime).String()
	url.SetHttpGetTimeoutSeconds(time.Duration(flagMaxTimeoutInSeconds) * time.Second)

	// create reports for all http urls
	var urlReports url.UrlReport
	if flagExecuteDryRun {
		urlReports = url.CreateDryReport(httpUrls)
	} else {
		urlReports = url.CustomizableCreateUrlReport(httpUrls, int(flagMaxParallelRequests))
	}
	urlReports.AddMetaData("initial_parsing_duration", parsing_duration)
	urlReports.AddMetaData("total_extracted_urls", fmt.Sprint(numberUniqueUrls))

	// cleanup reports if wanted
	if !flagUrlShowReachables {
		urlReports = urlReports.CleanupReachableUrls()
	}

	// creating report output
	var reportOutput string
	switch {
	case flagOutputJson:
		reportOutput = urlReports.Json()
	case flagOutputCsv:
		reportOutput = urlReports.FullString()
	default:
		fmt.Println("No output format was chosen, that should never happen. How did you do that?")
	}

	//  deciding where to write output to
	if flagOutputInFile != "" {
		err = url.WriteTo(flagOutputInFile, reportOutput)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(reportOutput)
	}
	fmt.Println(goodbyMsg)
	if !urlReports.AllReachable() {
		os.Exit(1)
	}
}
