package url

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUrlIsAvailable(t *testing.T) {

	t.Run("handeling 200er results in time", func(t *testing.T) {
		serverStatusCode := 200
		fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(serverStatusCode)
		}))
		defer fakeServer.Close()
		got := UrlIsAvailable(ExtractedUrl{fakeServer.URL, 1})
		want := UrlStatus{fakeServer.URL, true, http.StatusText(serverStatusCode), 1}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("handeling none 200er results in time", func(t *testing.T) {
		serverStatusCode := 404
		fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(serverStatusCode)
		}))
		defer fakeServer.Close()
		got := UrlIsAvailable(ExtractedUrl{fakeServer.URL, 2})
		want := UrlStatus{fakeServer.URL, false, http.StatusText(serverStatusCode), 2}

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

	})

	t.Run("server need more than 5ms to respond (timeout on dunction is configurable)", func(t *testing.T) {
		serverResponseCode := 200
		crawlerTimeout := time.Millisecond * 5
		fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(serverResponseCode)
		}))
		defer fakeServer.Close()
		SetHttpGetTimeoutSeconds(crawlerTimeout)
		got := UrlIsAvailable(ExtractedUrl{fakeServer.URL, 3})
		want := UrlStatus{fakeServer.URL, false, createTimeoutMessage(crawlerTimeout), 3}
		SetHttpGetTimeoutSeconds(DefaultHttpGetTimeout)
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
