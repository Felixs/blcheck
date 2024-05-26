package report

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Felixs/blcheck/pkg/url"
)

func TestCreateUrlReport(t *testing.T) {

	t.Run("Check with empty url slice", func(t *testing.T) {
		r := CreateUrlReport([]url.ExtractedUrl{})
		if len(r.UrlStatus) != 0 {
			t.Errorf("Expected report to be empty, got %d entries", len(r.UrlStatus))
		}
	})

	t.Run("Check with one entry", func(t *testing.T) {
		fakeServer := createDelayServerWithStatus(0*time.Second, 200)
		defer fakeServer.Close()
		inputUrls := []url.ExtractedUrl{{Url: fakeServer.URL, NumOccured: 1}}
		expectedUrlStatus := []url.UrlStatus{{Url: fakeServer.URL, IsReachable: true, StatusMessage: http.StatusText(200)}}
		r := CreateUrlReport(inputUrls)
		assertReport(t, r, expectedUrlStatus)
	})

	t.Run("Check with two entry", func(t *testing.T) {
		fakeServer := createDelayServerWithStatus(0*time.Second, 200)
		defer fakeServer.Close()
		fakeServer2 := createDelayServerWithStatus(0*time.Second, 404)
		defer fakeServer2.Close()
		inputUrls := []url.ExtractedUrl{
			{Url: fakeServer.URL, NumOccured: 1},
			{Url: fakeServer2.URL, NumOccured: 2},
		}
		expectedUrlStatus := []url.UrlStatus{
			{Url: fakeServer.URL, IsReachable: true, StatusMessage: ""},
			{Url: fakeServer2.URL, IsReachable: false, StatusMessage: ""},
		}
		r := CreateUrlReport(inputUrls)
		assertReport(t, r, expectedUrlStatus)
	})

	t.Run("runs checks in parallel", func(t *testing.T) {
		fakeServer := createDelayServerWithStatus(50*time.Millisecond, 200)
		defer fakeServer.Close()
		inputUrls := []url.ExtractedUrl{
			{Url: fakeServer.URL, NumOccured: 1},
			{Url: fakeServer.URL, NumOccured: 2},
			{Url: fakeServer.URL, NumOccured: 3},
		}
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
		inputUrls := []url.ExtractedUrl{
			{Url: fakeServer.URL, NumOccured: 1},
			{Url: fakeServer.URL, NumOccured: 2},
			{Url: fakeServer.URL, NumOccured: 3},
		}
		urlReport := CustomizableCreateUrlReport(inputUrls, 1)
		minRuntime := 150 * time.Millisecond
		if urlReport.Runtime < minRuntime {
			t.Errorf("Runtime was to fast with %v expected more than %v", urlReport.Runtime, minRuntime)
		}
	})
}

func TestAddMetaData(t *testing.T) {
	t.Run("Create with meta data", func(t *testing.T) {
		report := UrlReport{
			ExecutedAt: time.Now(),
			Runtime:    time.Second,
			MetaData:   map[string]string{"test": "testvalue"},
			UrlStatus:  []url.UrlStatus{},
		}
		want := map[string]string{"test": "testvalue"}

		if !reflect.DeepEqual(want, report.MetaData) {
			t.Errorf("got %v expected %v", report.MetaData, want)
		}
	})
	t.Run("Add MetaData afterwards", func(t *testing.T) {
		report := NewUrlReport(
			time.Now(),
			time.Second,
			[]url.UrlStatus{},
		)
		key := "test2"
		value := "new test value"
		want := map[string]string{key: value}
		report.AddMetaData(key, value)

		if !reflect.DeepEqual(want, report.MetaData) {
			t.Errorf("got %v expected %v", report.MetaData, want)
		}
	})

	t.Run("Overwrite MetaData", func(t *testing.T) {
		report := UrlReport{
			time.Now(),
			time.Second,
			map[string]string{"test2": "first set data"},
			[]url.UrlStatus{},
		}
		key := "test2"
		value := "new test value"
		want := map[string]string{key: value}
		report.AddMetaData(key, value)

		if !reflect.DeepEqual(want, report.MetaData) {
			t.Errorf("got %v expected %v", report.MetaData, want)
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
