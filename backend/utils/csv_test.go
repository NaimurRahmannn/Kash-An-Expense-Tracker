package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestEnsureFileExistsCreatesFileWithHeader(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "users.csv")
	header := []string{"id", "name", "email"}

	if err := EnsureFileExists(filePath, header); err != nil {
		t.Fatalf("expected file to be created: %v", err)
	}

	rows, err := ReadCSV(filePath)
	if err != nil {
		t.Fatalf("expected file to be readable: %v", err)
	}

	if !reflect.DeepEqual(rows, [][]string{header}) {
		t.Fatalf("expected header row %v, got %v", header, rows)
	}
}

func TestAppendCSVAppendsRow(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "users.csv")
	header := []string{"id", "name"}
	row := []string{"1", "John Doe"}

	if err := EnsureFileExists(filePath, header); err != nil {
		t.Fatalf("expected file to be created: %v", err)
	}

	if err := AppendCSV(filePath, row); err != nil {
		t.Fatalf("expected row to append: %v", err)
	}

	rows, err := ReadCSV(filePath)
	if err != nil {
		t.Fatalf("expected file to be readable: %v", err)
	}

	expected := [][]string{header, row}
	if !reflect.DeepEqual(rows, expected) {
		t.Fatalf("expected rows %v, got %v", expected, rows)
	}
}

func TestReadCSVReadsHeaderAndRow(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "users.csv")
	expected := [][]string{
		{"id", "name"},
		{"1", "John Doe"},
	}

	if err := WriteCSV(filePath, expected); err != nil {
		t.Fatalf("expected file to be written: %v", err)
	}

	rows, err := ReadCSV(filePath)
	if err != nil {
		t.Fatalf("expected file to be readable: %v", err)
	}

	if !reflect.DeepEqual(rows, expected) {
		t.Fatalf("expected rows %v, got %v", expected, rows)
	}
}

func TestWriteCSVRewritesFullContent(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "users.csv")
	original := [][]string{
		{"id", "name"},
		{"1", "John Doe"},
	}
	replacement := [][]string{
		{"id", "name"},
		{"2", "Jane Doe"},
	}

	if err := WriteCSV(filePath, original); err != nil {
		t.Fatalf("expected original content to be written: %v", err)
	}

	if err := WriteCSV(filePath, replacement); err != nil {
		t.Fatalf("expected replacement content to be written: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("expected file to be readable: %v", err)
	}

	if string(content) != "id,name\n2,Jane Doe\n" {
		t.Fatalf("expected replacement content only, got %q", string(content))
	}
}
