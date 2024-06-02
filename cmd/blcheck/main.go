/*
blcheck - A simple tool to check which links on your websites are broken.
*/
package main

import (
	"fmt"
	"os"
	"time"

	args "github.com/Felixs/blcheck/pkg/arguments" // handels flag parsing on init
	"github.com/Felixs/blcheck/pkg/url"
)

// Blcheck entry point.
func main() {
	args.Parse()
	parseStart := time.Now()
	httpUrls, err := extractURLs(args.URL)
	if err != nil {
		fmt.Printf("Failure to extract links from given URL. ERROR: %v", err)
		os.Exit(6)
	}
	parsingDuration := time.Since(parseStart)

	// create reports for all http urls
	urlReports := createUrlReport(httpUrls)
	urlReports.AddMetaData("initial_parsing_duration", parsingDuration.String())

	// creating report in desired output and format
	err = deliverReport(urlReports)
	if err != nil {
		fmt.Printf("Failure to deliver output. ERROR: %v", err)
		os.Exit(7)
	}
	fmt.Println(args.GoodbyMsg)

	// descide on exit code
	if !urlReports.AllReachable() {
		os.Exit(1)
	}
}

// Delivers UrlReport as desired format
func deliverReport(urlReports url.UrlReport) error {
	var reportOutput string
	var err error
	switch {
	case args.OutputAsJSON:
		reportOutput, err = urlReports.Json()
	case args.OutputAsCSV:
		reportOutput, err = urlReports.Csv(true)
	default:
		reportOutput = urlReports.FullString()
	}
	if err != nil {
		fmt.Println("Error in report output creation: " + err.Error())
	}

	if args.OutputInFile != "" {
		err := url.WriteTo(args.OutputInFile, reportOutput)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(reportOutput)
	}
	return nil
}

// Creates UrlReport for given urls
func createUrlReport(httpUrls []url.ExtractedUrl) url.UrlReport {
	var urlReports url.UrlReport

	switch args.ExecuteDryRun {
	case true:
		urlReports = url.CreateDryReport(httpUrls)
	default:
		urlReports = url.CustomizableCreateUrlReport(httpUrls, args.MaxParallelRequests)
	}

	urlReports.AddMetaData("total_extracted_urls", fmt.Sprint(len(httpUrls)))

	if !args.ShowReachables {
		urlReports = urlReports.CleanupReachableUrls()
	}
	return urlReports
}

// Reads url and extracts unique urls with count of occurences
func extractURLs(inputUrl string) ([]url.ExtractedUrl, error) {
	fmt.Println("Checking URL: ", inputUrl)

	body, err := url.GetBodyFromUrl(inputUrl)
	if err != nil {
		return nil, err
	}
	httpUrls := url.ExtractHttpUrls(body)

	// check for exclusion
	if args.RegexExclude != "" {
		httpUrls = url.FilterByExclude(httpUrls, args.RegexExclude)
	}
	// check for inclusion
	if args.RegexInclude != "" {
		httpUrls = url.FilterByInclude(httpUrls, args.RegexInclude)
	}
	return httpUrls, nil
}
