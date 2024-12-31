package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	if err != nil {
		return false
	}

	return true
}

func CopyFile(src, dst string) error {
	inFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, inFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	// Ensure the destination file is properly flushed
	err = outFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to flush destination file: %v", err)
	}

	return nil
}
