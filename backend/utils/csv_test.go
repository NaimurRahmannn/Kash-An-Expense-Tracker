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

func TestEnsureFileExistsDoesNotDuplicateHeader(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "users.csv")
	header := []string{"id", "name"}

	if err := EnsureFileExists(filePath, header); err != nil {
		t.Fatalf("expected file to be created: %v", err)
	}
	if err := EnsureFileExists(filePath, header); err != nil {
		t.Fatalf("expected existing file check to succeed: %v", err)
	}

	rows, err := ReadCSV(filePath)
	if err != nil {
		t.Fatalf("expected file to be readable: %v", err)
	}

	if !reflect.DeepEqual(rows, [][]string{header}) {
		t.Fatalf("expected one header row, got %v", rows)
	}
}

func TestEnsureFileExistsCreatesNestedParentDirectories(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "nested", "deeper", "users.csv")
	header := []string{"id", "name"}

	if err := EnsureFileExists(filePath, header); err != nil {
		t.Fatalf("expected nested file to be created: %v", err)
	}

	rows, err := ReadCSV(filePath)
	if err != nil {
		t.Fatalf("expected nested file to be readable: %v", err)
	}

	if !reflect.DeepEqual(rows, [][]string{header}) {
		t.Fatalf("expected header row %v, got %v", header, rows)
	}
}

func TestEnsureFileExistsReturnsErrorWhenParentPathIsFile(t *testing.T) {
	parentPath := filepath.Join(t.TempDir(), "parent")
	if err := os.WriteFile(parentPath, []byte("not a directory"), 0644); err != nil {
		t.Fatalf("expected parent file to be written: %v", err)
	}

	err := EnsureFileExists(filepath.Join(parentPath, "users.csv"), []string{"id", "name"})
	if err == nil {
		t.Fatal("expected parent path file to return an error")
	}
}

func TestEnsureFileExistsCreatesFileWithoutHeader(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "empty.csv")

	if err := EnsureFileExists(filePath, nil); err != nil {
		t.Fatalf("expected file to be created without header: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("expected file to be readable: %v", err)
	}

	if string(content) != "" {
		t.Fatalf("expected empty file, got %q", string(content))
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

func TestAppendCSVWorksWithBareFilename(t *testing.T) {
	t.Chdir(t.TempDir())

	row := []string{"1", "John Doe"}
	if err := AppendCSV("users.csv", row); err != nil {
		t.Fatalf("expected row to append to bare filename: %v", err)
	}

	rows, err := ReadCSV("users.csv")
	if err != nil {
		t.Fatalf("expected bare filename to be readable: %v", err)
	}

	if !reflect.DeepEqual(rows, [][]string{row}) {
		t.Fatalf("expected rows %v, got %v", [][]string{row}, rows)
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

func TestReadCSVReturnsErrorForMissingFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "missing.csv")

	if _, err := ReadCSV(filePath); err == nil {
		t.Fatal("expected missing file to return an error")
	}
}

func TestReadCSVHandlesEmptyFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "empty.csv")

	if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
		t.Fatalf("expected empty file to be written: %v", err)
	}

	rows, err := ReadCSV(filePath)
	if err != nil {
		t.Fatalf("expected empty file to be readable: %v", err)
	}

	if len(rows) != 0 {
		t.Fatalf("expected no rows, got %v", rows)
	}
}

func TestWriteCSVReturnsErrorWhenParentPathIsFile(t *testing.T) {
	parentPath := filepath.Join(t.TempDir(), "parent")
	if err := os.WriteFile(parentPath, []byte("not a directory"), 0644); err != nil {
		t.Fatalf("expected parent file to be written: %v", err)
	}

	err := WriteCSV(filepath.Join(parentPath, "users.csv"), [][]string{{"id", "name"}})
	if err == nil {
		t.Fatal("expected parent path file to return an error")
	}
}

func TestWriteCSVReturnsErrorWhenFilePathIsDirectory(t *testing.T) {
	dirPath := t.TempDir()

	err := WriteCSV(dirPath, [][]string{{"id", "name"}})
	if err == nil {
		t.Fatal("expected directory file path to return an error")
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

func TestAppendCSVReturnsErrorWhenParentPathIsFile(t *testing.T) {
	parentPath := filepath.Join(t.TempDir(), "parent")
	if err := os.WriteFile(parentPath, []byte("not a directory"), 0644); err != nil {
		t.Fatalf("expected parent file to be written: %v", err)
	}

	err := AppendCSV(filepath.Join(parentPath, "users.csv"), []string{"1", "John Doe"})
	if err == nil {
		t.Fatal("expected parent path file to return an error")
	}
}
