package url

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"
)

func TestIsUrlValid(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want bool
	}{
		{
			"valid https url",
			"https://www.google.de",
			true,
		}, {
			"valid http url",
			"https://www.google.de",
			true,
		}, {
			"valid with path",
			"https://www.google.de/gelbeseiten",
			true,
		}, {
			"invalid url",
			"wwwgooglede",
			false,
		}, {
			"missing protocol",
			"www.google.de",
			false,
		}, {
			"missing host name",
			"https://",
			false,
		}, {
			"empty string",
			"   ",
			false,
		}, {
			"no tld",
			"https://deinemama",
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUrlValid(tt.url)
			if got != tt.want {
				t.Errorf("got %v want %v on %s", got, tt.want, tt.url)
			}
		})
	}
}

func TestGetBodyFromUrl(t *testing.T) {
	t.Run("check body from a page that has 200 return", func(t *testing.T) {
		bodyData := "<html><body>hello</body></html>"
		fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte(bodyData))
		}))
		defer fakeServer.Close()

		url := fakeServer.URL
		got, _ := GetBodyFromUrl(url)
		want := bodyData

		if got != want {
			t.Errorf("got unexpected body %q from url %s", got, url)
		}
	})
	t.Run("check body from a page that has a none 200 return", func(t *testing.T) {
		fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
		}))
		defer fakeServer.Close()

		url := fakeServer.URL
		_, err := GetBodyFromUrl(url)

		if err == nil {
			t.Fatal("expected to get an error")
		}
	})
}

// infers https prefix if url does not start with a http protocol
func TestInferHttpsPrefix(t *testing.T) {

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			"adding https for missing protocol",
			"www.heise.de",
			"https://www.heise.de",
		}, {
			"check for existing https:// prefix",
			"https://www.heise.de",
			"https://www.heise.de",
		}, {
			"check for existing http:// prefix",
			"http://www.heise.de",
			"http://www.heise.de",
		}, {
			"check for empty string, seems stupid, but lets do it",
			"",
			"https://",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			InferHttpsPrefix(&tt.input)

			if tt.input != tt.want {
				t.Errorf("got %s want %s", tt.input, tt.want)
			}
		})
	}
}

func TestExtractHttpUrls(t *testing.T) {

	cases := []struct {
		name string
		body string
		want []string
	}{
		{
			"basic one url in body",
			`<html><body><a href="http://www.google.de">test</a></body></html>`,
			[]string{"http://www.google.de"},
		}, {
			"No urls found",
			`<html><body><a href="google.de">test</a></body></html>`,
			[]string{},
		}, {
			"No urls found",
			`<html><body><a href="google.de">https://heise.de http://www.google.de</a></body></html>`,
			[]string{"https://heise.de", "http://www.google.de"},
		}, {
			"Only non http urls",
			`<html><body><a href="mailto://google.de">file://heise.de xmpp://www.google.de</a></body></html>`,
			[]string{},
		}, {
			"Remove doubled urls",
			`<html><body><a href="http://www.google.de">http://www.google.de</a></body></html>`,
			[]string{"http://www.google.de"},
		}, {
			"Lowercase and uppercase have to be ignored, cast everything to lowercast",
			`<html><body><a href="http://www.GOOGLE.de">http://www.google.de</a></body></html>`,
			[]string{"http://www.google.de"},
		}, {
			"Cut ancor tags on links",
			`<html><body><a href="http://www.google.de/#very-good-link">http://www.google.de/</a></body></html>`,
			[]string{"http://www.google.de"},
		}, {
			"Remove tailing / on urls",
			`<html><body><a href="http://www.google.de/">hello</a></body></html>`,
			[]string{"http://www.google.de"},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractHttpUrls(tt.body)

			if len(got) != len(tt.want) {
				t.Fatalf("expected %v and %v to have the same length", tt.want, got)
			}
			for _, wantElement := range tt.want {
				if !slices.Contains(got, wantElement) {
					t.Errorf("expected %v to include want %s", got, wantElement)
				}
			}
		})
	}
}

func TestUrlIsAvailable(t *testing.T) {
	t.Run("server returns 200", func(t *testing.T) {
		fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer fakeServer.Close()

		got := UrlIsAvailable(fakeServer.URL)
		want := true

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("server returns 404", func(t *testing.T) {
		fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
		}))
		defer fakeServer.Close()

		got := UrlIsAvailable(fakeServer.URL)
		want := false

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("server need more than 5ms to respond (timeout on dunction is configurable)", func(t *testing.T) {
		fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(200)
		}))
		defer fakeServer.Close()

		got := ConfigurableUrlIsAvailable(fakeServer.URL, time.Millisecond*5)
		want := false

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
