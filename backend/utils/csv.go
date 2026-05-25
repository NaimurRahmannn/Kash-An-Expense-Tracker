package utils

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

// EnsureFileExists creates a CSV file with the provided header when it is missing.
func EnsureFileExists(filePath string, header []string) error {
	if _, err := os.Stat(filePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := ensureParentDir(filePath); err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	var writeErr error
	if len(header) > 0 {
		writeErr = writer.Write(header)
	}
	writer.Flush()

	if writeErr != nil {
		return writeErr
	}

	return writer.Error()
}

// ReadCSV reads all rows from a CSV file.
func ReadCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

// WriteCSV rewrites a CSV file with the provided rows.
func WriteCSV(filePath string, rows [][]string) error {
	if err := ensureParentDir(filePath); err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			writer.Flush()
			return err
		}
	}
	writer.Flush()

	return writer.Error()
}

// AppendCSV appends a single row to a CSV file.
func AppendCSV(filePath string, row []string) error {
	if err := ensureParentDir(filePath); err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writeErr := writer.Write(row)
	writer.Flush()

	if writeErr != nil {
		return writeErr
	}

	return writer.Error()
}

func ensureParentDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "" {
		return nil
	}

	return os.MkdirAll(dir, 0755)
}
