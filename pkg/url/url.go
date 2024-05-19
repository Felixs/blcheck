package url

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// checks if given url string seems valid
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

// tries to recieve
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
