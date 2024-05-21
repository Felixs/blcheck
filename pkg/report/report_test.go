package report

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Felixs/blcheck/pkg/url"
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
		expectedUrlStatus := []url.UrlStatus{{fakeServer.URL, true, http.StatusText(200)}}
		r := CreateUrlReport(inputUrls)
		assertReport(t, r, expectedUrlStatus)
	})

	t.Run("Check with two entry", func(t *testing.T) {
		fakeServer := createDelayServerWithStatus(0*time.Second, 200)
		defer fakeServer.Close()
		fakeServer2 := createDelayServerWithStatus(0*time.Second, 404)
		defer fakeServer2.Close()
		inputUrls := []string{fakeServer.URL, fakeServer2.URL}
		expectedUrlStatus := []url.UrlStatus{{fakeServer.URL, true, ""}, {fakeServer2.URL, false, ""}}
		r := CreateUrlReport(inputUrls)
		assertReport(t, r, expectedUrlStatus)
	})

	t.Run("runs checks in parallel", func(t *testing.T) {
		fakeServer := createDelayServerWithStatus(50*time.Millisecond, 200)
		defer fakeServer.Close()
		inputUrls := []string{fakeServer.URL, fakeServer.URL, fakeServer.URL}
		urlReport := CreateUrlReport(inputUrls)
		maxRuntime := 75 * time.Millisecond
		// flaky test, >75 is half of execution time in squence
		if urlReport.Runtime > maxRuntime {
			t.Errorf("Runtime was to slow with %v expected less than %v", urlReport.Runtime, maxRuntime)
		}
	})

	t.Run("runs checks in parallel with max number of workern", func(t *testing.T) {
		fakeServer := createDelayServerWithStatus(50*time.Millisecond, 200)
		defer fakeServer.Close()
		inputUrls := []string{fakeServer.URL, fakeServer.URL, fakeServer.URL}
		urlReport := CustomizableCreateUrlReport(inputUrls, 1)
		minRuntime := 150 * time.Millisecond
		if urlReport.Runtime < minRuntime {
			t.Errorf("Runtime was to fast with %v expected more than %v", urlReport.Runtime, minRuntime)
		}
	})
}

func assertReport(t *testing.T, r UrlReport, expectedUrlStatus []url.UrlStatus) {
	t.Helper()
	if len(r.UrlStatus) != len(expectedUrlStatus) {
		t.Fatalf("Checked urls not the same length, got %d want %d", len(r.UrlStatus), len(expectedUrlStatus))
	}
	// TODO: need a better solution to check if output is okay, this seems overkill
	lastFound := false
	for _, tt := range expectedUrlStatus {
		lastFound = false
		for _, bt := range r.UrlStatus {
			if tt.Url == bt.Url {
				if tt.IsReachable != bt.IsReachable {
					t.Errorf("expected %q to be IsReachable=%v", tt.Url, tt.IsReachable)
				}
				lastFound = true
				break
			}
		}
		if lastFound == false {
			t.Errorf("expected %q to be in %v", tt.Url, r)
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
