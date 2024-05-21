package url

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"mvdan.cc/xurls/v2"
)

const (
	// time to wait for an answer of webserver
	defaultHttpGetTimeout = 5 * time.Second
	// protocol prefixes
	prefixHttp  = "http://"
	prefixHttps = "https://"
)

// Information of a availability check on one webpage.
type UrlStatus struct {
	Url           string `json:"url"`
	IsReachable   bool   `json:"is_reachable"`
	StatusMessage string `json:"status_message"`
}

// String representation of a UrlStatus.
func (s UrlStatus) String() string {
	return fmt.Sprintf("%v\t%s\t%s", s.IsReachable, s.StatusMessage, s.Url)
}

// Checks if given url string seems to be valid.
func IsUrlValid(inputUrl string) (isValid bool) {
	urlData, err := url.Parse(inputUrl)
	if err != nil {
		return false
	} else if !urlData.IsAbs() {
		return false
	} else if urlData.Host == "" {
		return false
	} else if !strings.Contains(urlData.Host, ".") {
		return false
	}
	return true
}

// Tries to recieve a body with get request from url and returns it as string.
func GetBodyFromUrl(inputUrl string) (body string, err error) {
	resp, err := http.Get(inputUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("got status that is not okay: " + resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	body = string(bodyBytes)
	return body, nil
}

// Extracts any uniqe http(s) url that can be found from a given string.
func ExtractHttpUrls(body string) (hrefs []string) {
	precompiledUrlRegex := xurls.Strict()
	strictUrls := precompiledUrlRegex.FindAllString(body, -1)

	filteredUrls := filterNoneHttpUrls(strictUrls)
	slices.Sort(filteredUrls)
	return filteredUrls
}

// Filters down to all unique http urls
func filterNoneHttpUrls(strictUrls []string) []string {
	httpUrls := []string{}
	for _, httpUrl := range strictUrls {
		if strings.HasPrefix(httpUrl, "http") {
			httpUrl = strings.ToLower(httpUrl)
			httpUrl = strings.SplitN(httpUrl, "#", 2)[0]
			httpUrl = strings.TrimSuffix(httpUrl, "/")
			if !slices.Contains(httpUrls, httpUrl) {
				httpUrls = append(httpUrls, httpUrl)
			}
		}
	}
	return httpUrls
}

// Given string is checked for prefixing http(s) protocoll and gets added https if needed.
func InferHttpsPrefix(inputUrl *string) {
	if !strings.HasPrefix(*inputUrl, prefixHttps) && !strings.HasPrefix(*inputUrl, prefixHttp) {
		log.Println("infered https:// prefix, because given url did not have an protocol")
		*inputUrl = "https://" + *inputUrl
	}
}

// Trys a Get request on url and if status code = 200 and within timeout of 3 seconds returns true. Otherwise false.
func UrlIsAvailable(inputUrl string) (available UrlStatus) {
	return ConfigurableUrlIsAvailable(inputUrl, defaultHttpGetTimeout)
}

// Trys a Get request on url and if status code = 200 and within timeout returns true. Otherwise false.
func ConfigurableUrlIsAvailable(inputUrl string, timeout time.Duration) (available UrlStatus) {
	select {
	case r := <-checkUrl(inputUrl):
		return r
	case <-time.After(timeout):
		return UrlStatus{
			Url:           inputUrl,
			IsReachable:   false,
			StatusMessage: createTimeoutMessage(timeout),
		}
	}
}

// Formats a duration as status message
func createTimeoutMessage(timeout time.Duration) string {
	return fmt.Sprintf("Timed out after %v", timeout)
}

// Creates chan that handels url get returns.
func checkUrl(inputUrl string) chan UrlStatus {
	ch := make(chan UrlStatus, 1)
	go func() {
		isReachable := false
		statusMessage := "Unknown"
		//FIXME: maybe use http.Head(inputUrl)?
		resp, err := http.Get(inputUrl)
		if err != nil {
			statusMessage = err.Error()
		} else {
			statusMessage = http.StatusText(resp.StatusCode)
			isReachable = (resp.StatusCode == http.StatusOK)
		}

		ch <- UrlStatus{Url: inputUrl, IsReachable: isReachable, StatusMessage: statusMessage}
	}()
	return ch
}
