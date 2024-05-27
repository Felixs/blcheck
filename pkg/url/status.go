package url

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultHttpGetTimeout = 5 * time.Second
)

// time to wait for an answer of webserver
var HttpGetTimeout = DefaultHttpGetTimeout

// Convinience list of all files in UrlStatus
var urlStatusHeader = []string{"url", "is_reachable", "status_message", "content_length", "response_time", "num_occured"}

// Information of a availability check on one webpage.
type UrlStatus struct {
	Url           string        `json:"url"`
	IsReachable   bool          `json:"is_reachable"`
	StatusMessage string        `json:"status_message"`
	ContentLength int64         `json:"content_length"`
	ResponseTime  time.Duration `json:"response_time"`
	NumOccured    int           `json:"num_occured"`
}

// String representation of a UrlStatus.
func (s UrlStatus) String() string {
	return fmt.Sprintf("%s\t%v\t%s\t%d\t%s\t%d", s.Url, s.IsReachable, s.StatusMessage, s.ContentLength, s.ResponseTime, s.NumOccured)
}

// Well formed flieds header string for UrlStatus
func UrlStatusHeaderString() string {
	var sb strings.Builder
	for _, s := range urlStatusHeader {
		sb.WriteString(s)
		sb.WriteString("\t")
	}
	return sb.String()
}

// Creates UrlStatus from ExtractedUrl with additional information
func UrlStatusFromExtractedUrl(e ExtractedUrl, isReachable bool, statusMessage string, contentLength int64, responseTime time.Duration) UrlStatus {
	return UrlStatus{
		Url:           e.Url,
		IsReachable:   isReachable,
		StatusMessage: statusMessage,
		ContentLength: contentLength,
		ResponseTime:  responseTime,
		NumOccured:    e.NumOccured,
	}
}

// Overwrites module wide timeout for requests
func SetHttpGetTimeoutSeconds(timeout time.Duration) {
	HttpGetTimeout = timeout
}

// Trys a Get request on url and if status code = 200 and within timeout of HttpGetTimeout. Otherwise false.
func UrlIsAvailable(inputUrl ExtractedUrl) (available UrlStatus) {
	return ConfigurableUrlIsAvailable(inputUrl, HttpGetTimeout)
}

// Trys a Get request on url and if status code = 200 and within timeout returns true. Otherwise false.
func ConfigurableUrlIsAvailable(inputUrl ExtractedUrl, timeout time.Duration) (available UrlStatus) {
	select {
	case r := <-checkUrl(inputUrl):
		return r
	case <-time.After(timeout):
		return UrlStatus{
			Url:           inputUrl.Url,
			IsReachable:   false,
			StatusMessage: createTimeoutMessage(timeout),
			ContentLength: -1,
			ResponseTime:  timeout,
			NumOccured:    inputUrl.NumOccured,
		}
	}
}

// Formats a duration as status message
func createTimeoutMessage(timeout time.Duration) string {
	return fmt.Sprintf("Timed out after %v", timeout)
}

// Creates chan that handels url get returns.
func checkUrl(inputUrl ExtractedUrl) chan UrlStatus {
	ch := make(chan UrlStatus, 1)
	go func() {
		isReachable := false
		statusMessage := "Unknown"
		var contenLength int64
		responseTime := 0 * time.Millisecond
		getTimerStart := time.Now()
		resp, err := http.Head(inputUrl.Url)
		if err != nil {
			statusMessage = err.Error()
		} else {
			responseTime = time.Since(getTimerStart)
			statusMessage = http.StatusText(resp.StatusCode)
			isReachable = (resp.StatusCode == http.StatusOK)
			contenLength = resp.ContentLength

		}
		ch <- UrlStatusFromExtractedUrl(inputUrl, isReachable, statusMessage, contenLength, responseTime)

	}()
	return ch
}
