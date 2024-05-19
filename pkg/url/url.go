package url

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Checks if given url string seems to be valid
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

// Tries to recieve a body with get request from url and returns it as string
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

func ExtractHrefs(body string) (hrefs []string) {
	return []string{""}
}

const (
	prefixHttp  = "http://"
	prefixHttps = "https://"
)

// Given string is checked for prefixing http(s) protocoll and gets added https if needed
func InferHttpsPrefix(inputUrl *string) {
	if !strings.HasPrefix(*inputUrl, prefixHttps) && !strings.HasPrefix(*inputUrl, prefixHttp) {
		log.Println("infered https:// prefix, because given url did not have an protocol")
		*inputUrl = "https://" + *inputUrl
	}
}
