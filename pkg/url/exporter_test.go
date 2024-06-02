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

	got, err := report.Json()
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
	if err != nil {
		t.Fatal("did not expect an error")
	}
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
	got, err := convertToJsonStruct(report)
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
	if err != nil {
		t.Fatal("did not expect an error")
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCsv(t *testing.T) {
	t.Run("one line", func(t *testing.T) {
		r := UrlReport{
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
		got, err := r.Csv(false)
		want := `https://www.google.de,true,OK,1000,1s,12
`
		if err != nil {
			t.Fatal("did not expect to get an error")
		}

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("two line", func(t *testing.T) {
		r := UrlReport{
			UrlStatus: []UrlStatus{
				{
					Url:           "https://www.google.de",
					IsReachable:   true,
					StatusMessage: "OK",
					ContentLength: 1000,
					ResponseTime:  time.Second,
					NumOccured:    12,
				}, {
					Url:           "https://www.google2.de",
					IsReachable:   false,
					StatusMessage: "Not Found",
					ContentLength: -1,
					ResponseTime:  time.Minute,
					NumOccured:    99,
				},
			},
		}
		got, err := r.Csv(false)
		want := `https://www.google.de,true,OK,1000,1s,12
https://www.google2.de,false,Not Found,-1,1m0s,99
`
		if err != nil {
			t.Fatal("did not expect to get an error")
		}
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("add header", func(t *testing.T) {
		r := UrlReport{
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
		got, err := r.Csv(true)
		want := `url,is_reachable,status_message,content_length,response_time,num_occured
https://www.google.de,true,OK,1000,1s,12
`
		if err != nil {
			t.Fatal("did not expect to get an error")
		}

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("check quote comma", func(t *testing.T) {
		r := UrlReport{
			UrlStatus: []UrlStatus{
				{
					Url:           "https://www,google.de",
					IsReachable:   true,
					StatusMessage: "O,K",
					ContentLength: 1000,
					ResponseTime:  time.Second,
					NumOccured:    12,
				},
			},
		}
		got, err := r.Csv(false)
		want := `"https://www,google.de",true,"O,K",1000,1s,12
`
		if err != nil {
			t.Fatal("did not expect to get an error")
		}

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
