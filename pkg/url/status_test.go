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
			w.Write([]byte("test"))
		}))
		defer fakeServer.Close()
		got := UrlIsAvailable(ExtractedUrl{fakeServer.URL, 1})
		want := UrlStatus{fakeServer.URL, true, http.StatusText(serverStatusCode), 4, time.Second, 1}

		valid, message := assertUrlStatus(want, got)
		if !valid {
			t.Errorf("got %v, want %v, error %s", got, want, message)
		}
	})

	t.Run("handeling none 200er results in time", func(t *testing.T) {
		serverStatusCode := 404
		fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(serverStatusCode)
			w.Write([]byte("test"))
		}))
		defer fakeServer.Close()
		got := UrlIsAvailable(ExtractedUrl{fakeServer.URL, 2})
		want := UrlStatus{fakeServer.URL, false, http.StatusText(serverStatusCode), 4, time.Second, 2}

		valid, message := assertUrlStatus(want, got)
		if !valid {
			t.Errorf("got %v, want %v, error %s", got, want, message)
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
		want := UrlStatus{fakeServer.URL, false, createTimeoutMessage(crawlerTimeout), -1, time.Second, 3}
		SetHttpGetTimeoutSeconds(DefaultHttpGetTimeout)
		valid, message := assertUrlStatus(want, got)
		if !valid {
			t.Errorf("got %v, want %v, error %s", got, want, message)
		}
	})
}

func assertUrlStatus(want UrlStatus, got UrlStatus) (bool, string) {
	if want.Url != got.Url ||
		want.IsReachable != got.IsReachable ||
		want.NumOccured != got.NumOccured ||
		want.StatusMessage != got.StatusMessage ||
		want.ContentLength != got.ContentLength {
		return false, "similarity check failed"
	}
	if want.ResponseTime < got.ResponseTime {
		return false, "ceiling check failed"
	}
	return true, ""

}
