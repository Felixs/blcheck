package report

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Felixs/blcheck/pkg/url"
)

func TestJson(t *testing.T) {
	report := UrlReport{
		ExecutedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Runtime:    10 * time.Second,
		UrlStatus: []url.UrlStatus{
			{
				Url:           "https://www.google.de",
				IsReachable:   true,
				StatusMessage: "OK",
				NumOccured:    1,
			},
		},
	}

	got := report.Json()
	want := `{
    "executed_at": "2024-01-01T00:00:00Z",
    "runtime": "10s",
    "url_status": [{
        "url": "https://www.google.de",
        "is_reachable": true,
        "status_message": "OK",
        "num_occured": 1
    }]
}`
	want = strings.ReplaceAll(want, " ", "")
	want = strings.ReplaceAll(want, "\n", "")
	if got != want {
		t.Errorf("got %q expected %q", got, want)
	}

}

func TestConvertToJsonStruct(t *testing.T) {
	report := UrlReport{
		ExecutedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Runtime:    10 * time.Second,
		UrlStatus: []url.UrlStatus{
			{
				Url:           "https://www.google.de",
				IsReachable:   true,
				StatusMessage: "OK",
			},
		},
	}
	got := convertToJsonStruct(report)
	want := JsonUrlReport{
		ExecutedAt: "2024-01-01T00:00:00Z",
		Runtime:    "10s",
		UrlStatus: []url.UrlStatus{
			{
				Url:           "https://www.google.de",
				IsReachable:   true,
				StatusMessage: "OK",
			},
		},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v, want %v", got, want)
	}

}
