package report

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Felixs/blcheck/pkg/url"
)

// Information about a url check on a single webpage
type UrlReport struct {
	ExecutedAt time.Time     `json:"executed_at"`
	Runtime    time.Duration `json:"runtime"`
	UrlStatus  []UrlStatus   `json:"url_status"`
}

// String representation auf UrlReport
func (r UrlReport) String() string {
	return fmt.Sprintf("Started: %s , took: %s, urlcount: %d", r.ExecutedAt, r.Runtime, len(r.UrlStatus))
}

// String representation auf UrlReport with all UrlStatus
func (r UrlReport) FullString() string {
	var builder strings.Builder
	builder.WriteString(r.String() + "\n")
	for i, s := range r.UrlStatus {
		index := fmt.Sprintf("#%d", i+1)
		builder.WriteString(index + "\t" + s.String() + "\n")
	}
	return builder.String()
}

// Information of a availability check on one webpage
type UrlStatus struct {
	Url           string `json:"url"`
	IsReachable   bool   `json:"is_reachable"`
	StatusMessage string `json:"status_message"`
}

// String representation of a UrlStatus
func (s UrlStatus) String() string {
	return fmt.Sprintf("%v\t%s\t%s", s.IsReachable, s.StatusMessage, s.Url)
}

// Creates UrlReport from a list of given urls
func CreateUrlReport(urls []string) UrlReport {
	// TODO: refactor this
	start := time.Now()
	resultChan := make(chan UrlStatus)
	urlStatus := []UrlStatus{}

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for n := range resultChan {
			urlStatus = append(urlStatus, n)
		}
	}()

	var wg sync.WaitGroup
	for _, checkUrl := range urls {
		wg.Add(1)
		go func(inputUrl string) {
			defer wg.Done()
			status := url.UrlIsAvailable(inputUrl)
			message := "OK"
			if !status {
				message = "Not Found"
			}
			resultChan <- UrlStatus{inputUrl, status, message}
		}(checkUrl)
	}

	wg.Wait()
	close(resultChan)
	wg2.Wait()

	return UrlReport{
		start,
		time.Since(start),
		urlStatus,
	}
}
