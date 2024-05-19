package url

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

// TODO: next test
// func TestExtractHrefs(t *testing.T) {
// 	bodyData := `<html><body><a href="www.google.de">test</a></body></html>`
// 	want := []string{"www.google.de"}
// 	got := ExtractHrefs(bodyData)

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("got %v want %v", got, want)
// 	}
// }
