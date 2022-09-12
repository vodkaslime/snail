package utils

import (
	"errors"
	"os"
)

func ClearFile(p string) error {
	// Clear existing file
	fstat, err := os.Stat(p)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err == nil && fstat.IsDir() {
		return errors.New("found file path as a dir")
	}

	if !errors.Is(err, os.ErrNotExist) {
		return os.Remove(p)
	}
	return nil
}
