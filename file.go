package main

import (
	"fmt"
	"io"
	"os"
)

func readFile(filePath string) (string, error) {
	sugar.Debugw("reading file", "file_path", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		return "", nil
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("error reading content of file: %w", err)
	}
	sugar.Debugw("successfully read file", "file_path", filePath)
	return string(content), nil
}

func saveInFile(content, filePath string) error {
	sugar.Debugw("opening file for writing", "file_path", filePath)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	sugar.Debugw("saving content in file", "content", content, "file_path", filePath)
	_, err = fmt.Fprint(f, content)
	if err != nil {
		return fmt.Errorf("error writing content to file: %w", err)
	}
	sugar.Debugw("successfully written content into file", "content", content, "file_path", filePath)
	return nil
}
