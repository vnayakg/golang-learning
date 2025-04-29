package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestValidateFilePath_FileDoesNotExist(t *testing.T) {
	path := "this_file_does_not_exist.txt"
	err := validateFilePath(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expected := "this_file_does_not_exist.txt: no such file exist"
	if err.Error() != expected {
		t.Errorf("expected file not exist error, got: %v", err)
	}
}

func TestValidateFilePath_PathIsDirectory(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	err := validateFilePath(dir)
	expected := fmt.Sprintf("%v: is a directory", dir)
	if expected != err.Error() {
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

func TestRun_InvalidFlag(t *testing.T) {
	err := run([]string{"-invalid"})
	if err == nil {
		println(err)
		t.Errorf("expected error for invalid flag")
	}
}

func TestRun(t *testing.T) {
	t.Run("for valid line flag", func(t *testing.T) {
		tmpfile := createTempFile(t, "line1\nline2\nline3\n")
		defer os.Remove(tmpfile)

		output := captureOutput(func() {
			if err := run([]string{"-l", tmpfile}); err != nil {
				t.Fatalf("run failed: %v", err)
			}
		})

		expected := fmt.Sprintf("       3 %v\n",
			tmpfile)
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})

	t.Run("for valid word flag", func(t *testing.T) {
		tmpfile := createTempFile(t, "line1\nline2\nline3\n")
		defer os.Remove(tmpfile)

		output := captureOutput(func() {
			if err := run([]string{"-w", tmpfile}); err != nil {
				t.Fatalf("run failed: %v", err)
			}
		})

		expected := fmt.Sprintf("       3 %v\n",
			tmpfile)
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})

	t.Run("for valid character flag", func(t *testing.T) {
		tmpfile := createTempFile(t, "line1\nline2\nline3\n")
		defer os.Remove(tmpfile)

		output := captureOutput(func() {
			if err := run([]string{"-c", tmpfile}); err != nil {
				t.Fatalf("run failed: %v", err)
			}
		})

		expected := fmt.Sprintf("       18 %v\n",
			tmpfile)
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})

	t.Run("for missing flag", func(t *testing.T) {
		tmpfile := createTempFile(t, "line1\nline2\nline3\n")
		defer os.Remove(tmpfile)

		output := captureOutput(func() {
			if err := run([]string{tmpfile}); err != nil {
				t.Fatalf("run failed: %v", err)
			}
		})

		expected := fmt.Sprintf("       3        3       18 %v\n",
			tmpfile)
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})

	t.Run("multiple files with flags", func(t *testing.T) {
		firstFile := createTempFile(t, "line1\nline2\nline3\n")
		secondFile := createTempFile(t, "")
		defer os.Remove(firstFile)
		defer os.Remove(secondFile)

		originalProcessFiles := processFiles
		processFiles = func(paths []string) ([]*FileResult, error) {

			return []*FileResult{{firstFile, 3, 3, 18, nil},
				{secondFile, 0, 0, 0, nil}}, nil
		}

		output := captureOutput(func() {
			if err := run([]string{"-l", "-w", firstFile, secondFile}); err != nil {
				t.Fatalf("run failed: %v", err)
			}
		})
		expected := fmt.Sprintf("       3        3 %v\n       0        0 %v\n       3        3 total",
			firstFile, secondFile)
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			t.Errorf("expected %q, got %q", expected, output)
		}
		processFiles = originalProcessFiles
	})

	t.Run("multiple files without flags", func(t *testing.T) {
		firstFile := createTempFile(t, "line1\nline2\nline3\n")
		secondFile := createTempFile(t, "")
		defer os.Remove(firstFile)
		defer os.Remove(secondFile)

		originalProcessFiles := processFiles
		processFiles = func(paths []string) ([]*FileResult, error) {

			return []*FileResult{{firstFile, 3, 3, 18, nil},
				{secondFile, 0, 0, 0, nil}}, nil
		}

		output := captureOutput(func() {
			if err := run([]string{firstFile, secondFile}); err != nil {
				t.Fatalf("run failed: %v", err)
			}
		})
		expected := fmt.Sprintf("       3        3       18 %v\n       0        0        0 %v\n       3        3       18 total\n",
			firstFile, secondFile)
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			t.Errorf("expected %q, got %q", expected, output)
		}
		processFiles = originalProcessFiles
	})

	t.Run("read from stdin without flags", func(t *testing.T) {
		cleanup := mockStdin("line1\nline2\nline3\n")
		defer cleanup()

		output := captureOutput(func() {
			if err := run([]string{}); err != nil {
				t.Fatalf("run failed: %v", err)
			}
		})

		expected := "       3        3       18\n"
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})

	t.Run("read from stdin with flags", func(t *testing.T) {
		cleanup := mockStdin("line1\nline2\nline3\n")
		defer cleanup()

		output := captureOutput(func() {
			if err := run([]string{"-l", "-w"}); err != nil {
				t.Fatalf("run failed: %v", err)
			}
		})

		expected := "       3        3\n"
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})
}

func TestCountFile(t *testing.T) {
	tmpfile := createTempFile(t, "\tline1\nline2\nline3\n")
	file, _ := os.Open(tmpfile)
	defer os.Remove(tmpfile)
	defer file.Close()

	lines, words, chars, err := countFileItems(bufio.NewReader(file))

	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if lines != 3 {
		t.Errorf("expected lines %v, got %v", 3, lines)
	}
	if words != 3 {
		t.Errorf("expected words %v, got %v", 3, words)
	}
	if chars != 19 {
		t.Errorf("expected chars %v, got %v", 19, chars)
	}
}

func TestProcessFilesWithWorkerPool_SingleFile(t *testing.T) {
	content := "hello world\nthis is a test\n"
	path := createTempFile(t, content)
	defer os.Remove(path)

	results, err := processFiles([]string{path})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	res := results[0]
	expectedLines := 2
	expectedWords := 6
	expectedChars := len(content)
	if res.lines != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, res.lines)
	}
	if res.words != expectedWords {
		t.Errorf("Expected %d words, got %d", expectedWords, res.words)
	}
	if res.chars != expectedChars {
		t.Errorf("Expected %d chars, got %d", expectedChars, res.chars)
	}
}

func TestProcessFilesWithWorkerPool_MultipleFiles(t *testing.T) {
	content1 := "line one\nline two"
	content2 := "a b c d\ne f"
	path1 := createTempFile(t, content1)
	path2 := createTempFile(t, content2)
	defer os.Remove(path1)
	defer os.Remove(path2)

	results, err := processFiles([]string{path1, path2})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
}
func TestPrintResult_AllFlags(t *testing.T) {
	options := &Options{showLines: true, showWords: true, showChars: true}
	expected := fmt.Sprintf("%8d %8d %8d %s\n", 2, 5, 25, "test.txt")

	output := captureOutput(func() {
		printResult(2, 5, 25, "test.txt", options)
	})

	if strings.TrimSpace(output) != strings.TrimSpace(expected) {
		t.Errorf("Expected output:\n%q\nGot:\n%q", expected, output)
	}
}

func TestPrintResult_OnlyLines(t *testing.T) {
	options := &Options{showLines: true}
	expected := fmt.Sprintf("%8d %s\n", 10, "file.go")

	output := captureOutput(func() {
		printResult(10, 20, 30, "file.go", options)
	})

	if strings.TrimSpace(output) != strings.TrimSpace(expected) {
		t.Errorf("Expected output:\n%q\nGot:\n%q", expected, output)
	}
}

func TestPrintResult_OnlyWords(t *testing.T) {
	options := &Options{showWords: true}
	expected := fmt.Sprintf("%8d %s\n", 42, "abc.txt")

	output := captureOutput(func() {
		printResult(0, 42, 0, "abc.txt", options)
	})

	if strings.TrimSpace(output) != strings.TrimSpace(expected) {
		t.Errorf("Expected output:\n%q\nGot:\n%q", expected, output)
	}
}

func TestPrintResult_OnlyChars(t *testing.T) {
	options := &Options{showChars: true}
	expected := fmt.Sprintf("%8d %s\n", 100, "log.txt")

	output := captureOutput(func() {
		printResult(0, 0, 100, "log.txt", options)
	})

	if strings.TrimSpace(output) != strings.TrimSpace(expected) {
		t.Errorf("Expected output:\n%q\nGot:\n%q", expected, output)
	}
}

func mockStdin(input string) func() {
	file, _ := os.CreateTemp("", "stdin-mock")
	file.WriteString(input)
	file.Seek(0, 0)

	origStdin := os.Stdin
	os.Stdin = file

	return func() {
		os.Stdin = origStdin
		file.Close()
		os.Remove(file.Name())
	}
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = stdout
	buf.ReadFrom(r)
	return buf.String()
}
