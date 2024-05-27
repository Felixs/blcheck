package url

import "testing"

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
