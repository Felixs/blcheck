package url

import (
	"errors"
	"os"
	"path/filepath"
)

// Write given string in output file at filepath, overwrites any existing files.
func WriteTo(writePath, output string) error {
	absWritePath, err := filepath.Abs(writePath)
	if err != nil {
		return errors.New("Could not convert " + writePath + " to absolute filepath")
	}

	err = os.WriteFile(absWritePath, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}
