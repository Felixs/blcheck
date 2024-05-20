package report

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestCreateUrlReport(t *testing.T) {

	t.Run("Check with empty url slice", func(t *testing.T) {
		r := CreateUrlReport([]string{})
		if len(r.UrlStatus) != 0 {
			t.Errorf("Expected report to be empty, got %d entries", len(r.UrlStatus))
		}
	})

	t.Run("Check with one entry", func(t *testing.T) {
		fakeServer := createDelayServerWithStatus(0*time.Second, 200)
		defer fakeServer.Close()
		inputUrls := []string{fakeServer.URL}
		expectedUrlStatus := []UrlStatus{{fakeServer.URL, true, http.StatusText(200)}}
		r := CreateUrlReport(inputUrls)
		assertReport(t, r, inputUrls, expectedUrlStatus)
	})

	t.Run("Check with two entry", func(t *testing.T) {
		fakeServer := createDelayServerWithStatus(0*time.Second, 200)
		defer fakeServer.Close()
		fakeServer2 := createDelayServerWithStatus(0*time.Second, 404)
		defer fakeServer2.Close()
		inputUrls := []string{fakeServer.URL, fakeServer2.URL}
		expectedUrlStatus := []UrlStatus{{fakeServer.URL, true, http.StatusText(200)}, {fakeServer2.URL, false, http.StatusText(404)}}
		r := CreateUrlReport(inputUrls)
		assertReport(t, r, inputUrls, expectedUrlStatus)
	})

	t.Run("runs checks in parallel", func(t *testing.T) {
		fakeServer := createDelayServerWithStatus(50*time.Millisecond, 200)
		defer fakeServer.Close()
		inputUrls := []string{fakeServer.URL, fakeServer.URL, fakeServer.URL}
		start := time.Now()
		CreateUrlReport(inputUrls)
		if time.Since(start) > 52*time.Millisecond {
			t.Errorf("Processing took to long")
		}

	})
}

func assertReport(t *testing.T, r UrlReport, inputUrls []string, expectedUrlStatus []UrlStatus) {
	t.Helper()
	if len(r.UrlStatus) != len(inputUrls) {
		t.Fatalf("Checked urls not the same length, got %d want %d", len(r.UrlStatus), len(inputUrls))
	}
	// TODO: need a better solution to check if output is okay
	for _, tt := range expectedUrlStatus {
		for _, bt := range r.UrlStatus {
			if tt.Url == bt.Url {
				if !reflect.DeepEqual(tt, bt) {
					t.Errorf("missmatch for same url: got %v expecetd %v", bt, tt)
				}
				break
			}
			//t.Errorf("%q not found", tt.Url)
		}
	}

}

func createDelayServerWithStatus(delay time.Duration, status int) *httptest.Server {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(status)
	}))

	return fakeServer
}
