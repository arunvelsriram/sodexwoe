package utils

import (
	"os"
	"path/filepath"
)

func CreateFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}

	return os.Create(path)
}
