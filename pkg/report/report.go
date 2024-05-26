package report

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Felixs/blcheck/pkg/url"
)

// Max number of parallel routines to query webserver.
const MaxNumParallelQueries = 20

// Information about a url check on a single webpage.
type UrlReport struct {
	ExecutedAt time.Time         `json:"executed_at"`
	Runtime    time.Duration     `json:"runtime"`
	MetaData   map[string]string `json:"meta_data"`
	UrlStatus  []url.UrlStatus   `json:"url_status"`
}

// Convinience constructor
func NewUrlReport(executedAt time.Time, runtime time.Duration, urlStatus []url.UrlStatus) UrlReport {
	return UrlReport{
		ExecutedAt: executedAt,
		Runtime:    runtime,
		MetaData:   map[string]string{},
		UrlStatus:  urlStatus,
	}

}

// String representation auf UrlReport.
func (r UrlReport) String() string {
	return fmt.Sprintf("Started: %s , took: %s, urlcount: %d", r.ExecutedAt, r.Runtime, len(r.UrlStatus))
}

// String representation auf UrlReport with all UrlStatus.
func (r UrlReport) FullString() string {
	var builder strings.Builder
	builder.WriteString(r.String() + "\n")
	if len(r.MetaData) > 0 {
		builder.WriteString("Meta information:\n")
		for k, v := range r.MetaData {
			builder.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}
	for i, s := range r.UrlStatus {
		index := fmt.Sprintf("#%d", i+1)
		builder.WriteString(index + "\t" + s.String() + "\n")
	}
	return builder.String()
}

// Add a key-value meta date for UrlReport, does overwrite exisiting MetaData keys
func (r UrlReport) AddMetaData(key, value string) {
	r.MetaData[key] = value
}

// Creates UrlReport from a list of given urls.
func CreateUrlReport(urls []url.ExtractedUrl) UrlReport {
	return CustomizableCreateUrlReport(urls, MaxNumParallelQueries)
}

// Creates UrlReport from a list of given urls with max. of parallel request routines.
func CustomizableCreateUrlReport(urls []url.ExtractedUrl, maxRoutines int) UrlReport {
	start := time.Now()
	inputChan := make(chan url.ExtractedUrl)
	resultChan := make(chan url.UrlStatus)
	urlStatus := []url.UrlStatus{}

	// go routine that reads from result chan
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go gatherResults(&wg2, resultChan, &urlStatus)

	// routine that reads from input chan and writes to result chan
	var wg sync.WaitGroup
	for i := 0; i < maxRoutines; i++ {
		wg.Add(1)
		go checkUrlHandler(inputChan, resultChan, &wg)
	}

	// write input to input chan
	for _, inputUrl := range urls {
		inputChan <- inputUrl
	}
	close(inputChan)
	wg.Wait()
	close(resultChan)
	wg2.Wait()

	return NewUrlReport(
		start,
		time.Since(start),
		urlStatus,
	)
}

// Go routine to gather UrlStatus from result channel.
func gatherResults(wg2 *sync.WaitGroup, resultChan chan url.UrlStatus, urlStatus *[]url.UrlStatus) {
	defer wg2.Done()
	for result := range resultChan {
		*urlStatus = append(*urlStatus, result)
	}
}

// Go routine to get run url string by UrlIsAvailable to get UrlStatus.
func checkUrlHandler(inputChan chan url.ExtractedUrl, resultChan chan url.UrlStatus, wg *sync.WaitGroup) {
	defer wg.Done()
	for inputUrl := range inputChan {
		status := url.UrlIsAvailable(inputUrl)
		resultChan <- status
	}
}
