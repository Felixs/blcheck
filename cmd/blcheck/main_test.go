package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Felixs/blcheck/pkg/constants"
)

const (
	baseBinName = "blcheck"
	baseBinPath = "test/"
	binPath     = baseBinPath + baseBinName
)

func TestBasicPageWithoutLinks(t *testing.T) {

	// build application
	binName, err := buildBinary()
	if err != nil {
		t.Fatalf("failed to build: %v", err)
	}
	// Start a web server to serve static HTML content\
	// for later tests
	// go startTestServer()

	cases := []struct {
		name          string
		arguments     string
		wantStatus    int
		wantOutputEnd string
	}{
		{
			name:          "Run without arguments",
			arguments:     "",
			wantStatus:    constants.ExitMissingParameter,
			wantOutputEnd: "ERROR:URL is required\n",
		}, {
			name:          "Run again down server",
			arguments:     "http://localhost:1337/index.html",
			wantStatus:    constants.ExitUrlNotReachable,
			wantOutputEnd: "ERROR: Failure to extract links from given URL.\n",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			output, statusCode, _ := execBlcheck(binName, tt.arguments)
			if statusCode != tt.wantStatus {
				t.Fatalf("Expected error code %d, got %d", tt.wantStatus, statusCode)
			}
			if !strings.HasSuffix(output, tt.wantOutputEnd) {
				t.Fatalf("Missing error message")
			}
		})
	}
}

// Starts execution of compiled programm with given arguments
func execBlcheck(binName, arguments string) (string, int, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", -1, err
	}

	cmdPath := filepath.Join(dir, binName)

	var cmd *exec.Cmd
	if len(arguments) == 0 {
		cmd = exec.Command(cmdPath)
	} else {
		cmd = exec.Command(cmdPath, arguments)
	}

	out, _ := cmd.Output()

	return string(out), cmd.ProcessState.ExitCode(), nil

}

// Builds current version of program in tmp folder
func buildBinary() (string, error) {
	build := exec.Command("go", "build", "-o", binPath)

	if err := build.Run(); err != nil {
		return "", fmt.Errorf("cannot build tool %s: %s", binPath, err)
	}
	return binPath, nil
}

// Serves test pages to execute blcheck against
// func startTestServer() error {
// 	http.Handle("/", http.FileServer(http.Dir("../../tests/page")))
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
