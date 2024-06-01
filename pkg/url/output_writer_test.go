package url_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Felixs/blcheck/pkg/url"
)

const testWriteFolder = "../../tests/"

func TestWriteTo(t *testing.T) {
	t.Run("Check simple write", func(t *testing.T) {
		filename := testWriteFolder + "outout.txt"
		writeContent := "file content"
		err := url.WriteTo(filename, writeContent)

		if err != nil {
			t.Errorf("Got unexpected error in write file, %v", err.Error())
		}
		assertWrittenFile(t, filename, writeContent)
	})

}

func assertWrittenFile(t *testing.T, fielpath, wantedContent string) {
	t.Helper()
	path, _ := filepath.Abs(fielpath)
	fileContent, err := os.ReadFile(path)

	if err != nil {
		t.Fatalf("got unexpected error %q", err.Error())
	}
	got := string(fileContent)

	if got != wantedContent {
		t.Errorf("got %q wanted %q", got, wantedContent)
	}
}
