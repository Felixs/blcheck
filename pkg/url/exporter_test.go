package url

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestJson(t *testing.T) {
	report := UrlReport{
		ExecutedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Runtime:    10 * time.Second,
		UrlStatus: []UrlStatus{
			{
				Url:           "https://www.google.de",
				IsReachable:   true,
				StatusMessage: "OK",
				ContentLength: 1000,
				ResponseTime:  5 * time.Second,
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
        "content_length": 1000,
        "response_time": "5s",
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
		UrlStatus: []UrlStatus{
			{
				Url:           "https://www.google.de",
				IsReachable:   true,
				StatusMessage: "OK",
				ContentLength: 1000,
				ResponseTime:  time.Second,
				NumOccured:    12,
			},
		},
	}
	got := convertToJsonStruct(report)
	want := JsonUrlReport{
		ExecutedAt: "2024-01-01T00:00:00Z",
		Runtime:    "10s",
		UrlStatus: []JsonUrlStatus{
			{
				Url:           "https://www.google.de",
				IsReachable:   true,
				StatusMessage: "OK",
				ContentLength: 1000,
				ResponseTime:  "1s",
				NumOccured:    12,
			},
		},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v, want %v", got, want)
	}
}
