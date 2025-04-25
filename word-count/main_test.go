package main

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestCountLines(t *testing.T) {
	t.Run("file with multiple lines", func(t *testing.T) {
		content := "line1\nline2\nline3\n"
		tmpFile := createTempFile(t, content)
		defer os.Remove(tmpFile)

		count, err := countLines(tmpFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 3 {
			t.Errorf("expected 3 lines, got %d", count)
		}
	})

	t.Run("empty file", func(t *testing.T) {
		tmpFile := createTempFile(t, "")
		defer os.Remove(tmpFile)

		count, err := countLines(tmpFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 0 {
			t.Errorf("expected 0 lines, got %d", count)
		}
	})
}

func TestValidateFilePath_FileDoesNotExist(t *testing.T) {
	path := "this_file_does_not_exist.txt"
	err := validateFilePath(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected file not exist error, got: %v", err)
	}
}

func TestValidateFilePath_PermissionDenied(t *testing.T) {
	tmpFile := createTempFile(t, "some content\n")
	defer os.Remove(tmpFile)

	err := os.Chmod(tmpFile, 0200)
	if err != nil {
		t.Fatalf("failed to chmod: %v", err)
	}
	defer os.Chmod(tmpFile, 0600)

	err = validateFilePath(tmpFile)
	if err == nil {
		t.Fatal("expected permission error, got nil")
	}
	if !errors.Is(err, os.ErrPermission) {
		t.Errorf("expected permission error, got: %v", err)
	}
}

func TestValidateFilePath_PathIsDirectory(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	err := validateFilePath(dir)
	if !strings.Contains(err.Error(), "path is a directory") {
		t.Errorf("expected directory error, got: %v", err)
	}
}

func TestValidateFilePath_ValidFile(t *testing.T) {
	file := createTempFile(t, "this is valid file")
	defer os.Remove(file)

	err := validateFilePath(file)
	if err != nil {
		t.Errorf("expected no error for valid file, got: %v", err)
	}
}

func createTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "testDir")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	return dir
}

func createTempFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}

func TestRun_MissingArgs(t *testing.T) {
	err := run([]string{})
	if err == nil || err.Error() != "please provide a file path" {
		t.Errorf("expected file path error, got %v", err)
	}
}

func TestRun_InvalidFlag(t *testing.T) {
	err := run([]string{"-invalid"})
	if err == nil {
		println(err)
		t.Errorf("expected error for invalid flag")
	}
}

func TestRun_LineCountSuccess(t *testing.T) {
	t.Run("for valid line flag", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "example.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())
		tmpfile.WriteString("line1\nline2\nline3\n")
		tmpfile.Close()

		err = run([]string{"-l", tmpfile.Name()})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("for valid word flag", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "example.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())
		tmpfile.WriteString("line1\nline2\nline3\n")
		tmpfile.Close()

		err = run([]string{"-w", tmpfile.Name()})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func TestCountWords(t *testing.T) {
	t.Run("file with multiple words", func(t *testing.T) {
		content := "line1 line2\nline3\n"
		tmpFile := createTempFile(t, content)
		defer os.Remove(tmpFile)

		count, err := countWords(tmpFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 3 {
			t.Errorf("expected 3 words, got %d", count)
		}
	})

	t.Run("empty file", func(t *testing.T) {
		tmpFile := createTempFile(t, "")
		defer os.Remove(tmpFile)

		count, err := countWords(tmpFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 0 {
			t.Errorf("expected 0 words, got %d", count)
		}
	})
}
