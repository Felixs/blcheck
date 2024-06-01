/*
blcheck - A simple tool to check which links on your websites are broken.
*/
package main

import (
	"fmt"
	"os"
	"time"

	args "github.com/Felixs/blcheck/pkg/arguments"
	"github.com/Felixs/blcheck/pkg/url"
)

// Blcheck entry point.
func main() {
	processUrl(args.URL)
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
	if args.RegexExclude != "" {
		httpUrls = url.FilterByExclude(httpUrls, args.RegexExclude)
	}
	// check for inclusion
	if args.RegexInclude != "" {
		httpUrls = url.FilterByInclude(httpUrls, args.RegexInclude)
	}

	parsing_duration := time.Since(processStarttime).String()
	url.SetHttpGetTimeoutSeconds(time.Duration(args.MaxTimeoutInSeconds) * time.Second)

	// create reports for all http urls
	var urlReports url.UrlReport
	if args.ExecuteDryRun {
		urlReports = url.CreateDryReport(httpUrls)
	} else {
		urlReports = url.CustomizableCreateUrlReport(httpUrls, args.MaxParallelRequests)
	}
	urlReports.AddMetaData("initial_parsing_duration", parsing_duration)
	urlReports.AddMetaData("total_extracted_urls", fmt.Sprint(numberUniqueUrls))

	// cleanup reports if wanted
	if !args.ShowReachables {
		urlReports = urlReports.CleanupReachableUrls()
	}

	// creating report output
	var reportOutput string
	switch {
	case args.OutputAsJSON:
		reportOutput = urlReports.Json()
	case args.OutputAsCSV:
		reportOutput = fmt.Sprintln("CSV format not impelented jet")
	default:
		reportOutput = urlReports.FullString()
	}

	//  deciding where to write output to
	if args.OutputInFile != "" {
		err = url.WriteTo(args.OutputInFile, reportOutput)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(reportOutput)
	}
	fmt.Println(args.GoodbyMsg)
	if !urlReports.AllReachable() {
		os.Exit(1)
	}
}
