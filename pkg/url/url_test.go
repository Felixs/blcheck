package url

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"testing"
)

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
	t.Run("Check basic unique url extraction", func(t *testing.T) {
		cases := []struct {
			name string
			body string
			want []ExtractedUrl
		}{
			{
				"basic one url in body",
				`<html><body><a href="http://www.google.de">test</a></body></html>`,
				[]ExtractedUrl{{"http://www.google.de", 1}},
			}, {
				"No urls found",
				`<html><body><a href="google.de">test</a></body></html>`,
				[]ExtractedUrl{},
			}, {
				"two unique urls found",
				`<html><body><a href="google.de">https://heise.de http://www.google.de</a></body></html>`,
				[]ExtractedUrl{
					{"https://heise.de", 1},
					{"http://www.google.de", 1},
				},
			}, {
				"Only non http urls",
				`<html><body><a href="mailto://google.de">file://heise.de xmpp://www.google.de</a></body></html>`,
				[]ExtractedUrl{},
			}, {
				"Remove doubled urls",
				`<html><body><a href="http://www.google.de">http://www.google.de</a></body></html>`,
				[]ExtractedUrl{{"http://www.google.de", 2}},
			}, {
				"Lowercase and uppercase have to be ignored, cast everything to lowercast",
				`<html><body><a href="http://www.GOOGLE.de">http://www.google.de</a></body></html>`,
				[]ExtractedUrl{{"http://www.google.de", 2}},
			}, {
				"Cut ancor tags on links",
				`<html><body><a href="http://www.google.de/#very-good-link">http://www.google.de/</a></body></html>`,
				[]ExtractedUrl{{"http://www.google.de", 2}},
			}, {
				"Remove tailing / on urls",
				`<html><body><a href="http://www.google.de/">hello</a></body></html>`,
				[]ExtractedUrl{{"http://www.google.de", 1}},
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
						t.Errorf("expected %v to include want %v", got, wantElement)
					}
				}
			})
		}

	})

}

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
		}, {
			"localhost without port",
			"http://localhost/index.html",
			true,
		},
		{
			"localhost with port",
			"http://localhost:1337/index.html",
			true,
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

func TestFilterByExclude(t *testing.T) {

	t.Run("exclude on an empty list", func(t *testing.T) {
		input := []ExtractedUrl{}
		got := FilterByExclude(input, "google")
		want := []ExtractedUrl{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("exclude on a list with one entry that needs to excludes", func(t *testing.T) {
		input := []ExtractedUrl{{Url: "www.google.de", NumOccured: 99}}
		got := FilterByExclude(input, "google")
		want := []ExtractedUrl{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("exclude on a list with one entry that needs to stay after exclude", func(t *testing.T) {
		input := []ExtractedUrl{{Url: "www.google.de", NumOccured: 99}}
		got := FilterByExclude(input, "heise")
		want := []ExtractedUrl{{Url: "www.google.de", NumOccured: 99}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("exclude on a list with two entry that needs to excludes", func(t *testing.T) {
		input := []ExtractedUrl{{Url: "http://www.google.de", NumOccured: 99}, {Url: "http://www.google.de/help", NumOccured: 99}}
		got := FilterByExclude(input, "google")
		want := []ExtractedUrl{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("exclude on a list with two entries one needs to stay after exclude", func(t *testing.T) {
		input := []ExtractedUrl{{Url: "www.google.de", NumOccured: 99}, {Url: "www.heise.de", NumOccured: 9}}
		got := FilterByExclude(input, "heise")
		want := []ExtractedUrl{{Url: "www.google.de", NumOccured: 99}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestFilterByInclude(t *testing.T) {

	t.Run("include on an empty list", func(t *testing.T) {
		input := []ExtractedUrl{}
		got := FilterByInclude(input, "google")
		want := []ExtractedUrl{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("include on a list with one entry that needs to be excluded", func(t *testing.T) {
		input := []ExtractedUrl{{Url: "www.google.de", NumOccured: 99}}
		got := FilterByInclude(input, "heise")
		want := []ExtractedUrl{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("include on a list with one entry that needs to stay after inclusion", func(t *testing.T) {
		input := []ExtractedUrl{{Url: "www.google.de", NumOccured: 99}}
		got := FilterByInclude(input, "google")
		want := []ExtractedUrl{{Url: "www.google.de", NumOccured: 99}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("include on a list with two entry that needs to be included", func(t *testing.T) {
		input := []ExtractedUrl{{Url: "http://www.google.de", NumOccured: 99}, {Url: "http://www.google.de/help", NumOccured: 99}}
		got := FilterByInclude(input, "google")
		want := []ExtractedUrl{{Url: "http://www.google.de", NumOccured: 99}, {Url: "http://www.google.de/help", NumOccured: 99}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("include on a list with two entries one needs to stay", func(t *testing.T) {
		input := []ExtractedUrl{{Url: "www.google.de", NumOccured: 99}, {Url: "www.heise.de", NumOccured: 9}}
		got := FilterByInclude(input, "google")
		want := []ExtractedUrl{{Url: "www.google.de", NumOccured: 99}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
