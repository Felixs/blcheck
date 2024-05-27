package url

import (
	"fmt"
	"net/http"
	"time"
)

// Information of a availability check on one webpage.
type UrlStatus struct {
	Url           string `json:"url"`
	IsReachable   bool   `json:"is_reachable"`
	StatusMessage string `json:"status_message"`
	NumOccured    int    `json:"num_occured"`
}

// String representation of a UrlStatus.
func (s UrlStatus) String() string {
	return fmt.Sprintf("%v\t%s\t%s\t%d", s.IsReachable, s.StatusMessage, s.Url, s.NumOccured)
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
		resp, err := http.Get(inputUrl.Url)
		if err != nil {
			statusMessage = err.Error()
		} else {
			statusMessage = http.StatusText(resp.StatusCode)
			isReachable = (resp.StatusCode == http.StatusOK)
		}
		ch <- inputUrl.ToUrlStatus(statusMessage, isReachable)

	}()
	return ch
}
