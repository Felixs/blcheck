package arguments

import (
	"testing"
)

func TestCheckArguments(t *testing.T) {
	URL = "at"
	err := checkUrlParameter(URL)
	if err == nil {
		t.Errorf("Expected an failure")
	}
}
