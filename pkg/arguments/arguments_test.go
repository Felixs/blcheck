package arguments

import (
	"testing"
)

func TestCheckArguments(t *testing.T) {
	t.Run("url string must contain at least 4 characters, 1 domain, 1 dot, 2 tld", func(t *testing.T) {
		URL = "ata"
		err := checkUrlParameter(URL)
		if err == nil {
			t.Errorf("Expected an failure")
		}
	})

	t.Run("url string must contain at least 1 domain", func(t *testing.T) {
		URL = "asdfasdfasfdasdf"
		err := checkUrlParameter(URL)
		if err == nil {
			t.Errorf("Expected an failure")
		}
	})

	t.Run("max one output format can be defined, trying one", func(t *testing.T) {
		formats := []bool{true}
		err := checkOutputFormats(formats)
		if err != nil {
			t.Errorf("No error expected, got %v", err)
		}
	})

	t.Run("max one output format can be defined, trying two one true", func(t *testing.T) {
		formats := []bool{true, false}
		err := checkOutputFormats(formats)
		if err != nil {
			t.Errorf("No error expected, got %v", err)
		}
	})

	t.Run("max one output format can be defined, trying two true", func(t *testing.T) {
		formats := []bool{true, true}
		err := checkOutputFormats(formats)
		if err == nil {
			t.Errorf("Expected error")
		}
	})

	t.Run("check functions need positiv integer input", func(t *testing.T) {
		cases := []struct {
			name   string
			method func(int) error
		}{
			{
				"MaxParallelRequests",
				checkMaxParallelRequests,
			}, {
				"MaxTimeoutInSeconds",
				checkMaxTimeoutInSeconds,
			},
		}
		for _, tt := range cases {
			t.Run("number of"+tt.name+"needs to be not 0", func(t *testing.T) {
				input := 0
				err := tt.method(input)
				if err == nil {
					t.Errorf("Expected error")
				}
			})

			t.Run("number of"+tt.name+"needs to be not negativ", func(t *testing.T) {
				input := -1
				err := tt.method(input)
				if err == nil {
					t.Errorf("Expected error")
				}
			})

			t.Run("number of"+tt.name+"needs to be positiv", func(t *testing.T) {
				input := 1
				err := tt.method(input)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			})
		}
	})

}
