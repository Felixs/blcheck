package url

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

const (
	// protocol prefixes
	prefixHttp  = "http://"
	prefixHttps = "https://"
)

// Contains parsing info about url that is to check
type ExtractedUrl struct {
	Url        string
	NumOccured int
}

// Convertion of ExtractedUrl to UrlStatus
func (e ExtractedUrl) ToUrlStatus(statusMessage string, isReachable bool) UrlStatus {
	return UrlStatus{
		Url:           e.Url,
		IsReachable:   isReachable,
		StatusMessage: statusMessage,
		NumOccured:    e.NumOccured,
	}
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
	// Get request to page
	resp, err := http.Get(inputUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Check return status code, if not 200 return error
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("got status that is not okay: " + resp.Status)
	}
	// Read all data from request body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	body = string(bodyBytes)
	return body, nil
}

// Extracts any uniqe http(s) url that can be found from a given string.
func ExtractHttpUrls(body string) (hrefs []ExtractedUrl) {
	precompiledUrlRegex := xurls.Strict()
	strictUrls := precompiledUrlRegex.FindAllString(body, -1)

	filteredUrls := filterNoneHttpUrls(strictUrls)
	return filteredUrls
}

// Filters down to all unique http urls
func filterNoneHttpUrls(strictUrls []string) []ExtractedUrl {
	httpUrls := make(map[string]int)
	for _, httpUrl := range strictUrls {
		if strings.HasPrefix(httpUrl, "http") {
			// convert all chars to lowercase, easy comparision
			httpUrl = strings.ToLower(httpUrl)
			// remove ancor sufixes, because the only point to a part of a single website
			httpUrl = strings.SplitN(httpUrl, "#", 2)[0]
			// relove following '/' because they are not needed
			httpUrl = strings.TrimSuffix(httpUrl, "/")
			httpUrls[httpUrl] += 1
		}
	}
	extractedUrls := []ExtractedUrl{}
	for k, v := range httpUrls {
		extractedUrls = append(extractedUrls, ExtractedUrl{k, v})
	}

	return extractedUrls
}

// Given string is checked for prefixing http(s) protocoll and gets added https if needed.
func InferHttpsPrefix(inputUrl *string) {
	if !strings.HasPrefix(*inputUrl, prefixHttps) && !strings.HasPrefix(*inputUrl, prefixHttp) {
		log.Println("infered https:// prefix, because given url did not have an protocol")
		*inputUrl = "https://" + *inputUrl
	}
}
