package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func validateFilePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%v: no such file exist", path)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("%v: permission denied", path)
		}
	}
	if info.IsDir() {
		return fmt.Errorf("%v: is a directory", path)
	}
	return nil
}

func countLines(path string) (int, error) {
	return countWithScanner(path, bufio.ScanLines)
}

func countWords(path string) (int, error) {
	return countWithScanner(path, bufio.ScanWords)
}

func countCharacters(path string) (int, error) {
	return countWithScanner(path, bufio.ScanRunes)
}

func countWithScanner(path string, split bufio.SplitFunc) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsPermission(err) {
			return 0, fmt.Errorf("%v: permission denied", path)
		}
		return 0, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(split)

	count := 0
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error scanning file: %w", err)
	}
	return count, nil
}

type FileResult struct {
	path                string
	lines, words, chars int
	err                 error
}

func run(args []string) error {
	flags := flag.NewFlagSet("wc", flag.ContinueOnError)
	showLines := flags.Bool("l", false, "Count lines")
	showWords := flags.Bool("w", false, "Count words")
	showChars := flags.Bool("c", false, "Count characters")
	if err := flags.Parse(args); err != nil {
		return err
	}
	showAll := !*showLines && !*showWords && !*showChars

	filepaths := flags.Args()
	if len(filepaths) == 0 {
		return fmt.Errorf("please provide a file path")
	}

	totalLines, totalWords, totalChars := 0, 0, 0
	fileResults := make([]FileResult, len(filepaths))

	for index, filepath := range filepaths {
		if err := validateFilePath(filepath); err != nil {
			return err
		}
		fileResults[index].path = filepath
		lines, words, chars := 0, 0, 0
		var err error

		if showAll || *showLines {
			lines, err = countLines(filepath)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			totalLines += lines
			fileResults[index].lines = lines
		}
		if showAll || *showWords {
			words, err = countWords(filepath)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			totalWords += words
			fileResults[index].words = words
		}
		if showAll || *showChars {
			chars, err = countCharacters(filepath)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			totalChars += totalChars
			fileResults[index].chars = chars
		}
	}

	for _, result := range fileResults {
		if result.err != nil {
			fmt.Fprintln(os.Stderr, result.err)
			continue
		}
		if showAll || *showLines {
			fmt.Printf("%8d ", result.lines)
		}
		if showAll || *showWords {
			fmt.Printf("%8d ", result.words)
		}
		if showAll || *showChars {
			fmt.Printf("%8d ", result.chars)
		}
		fmt.Printf("%s\n", result.path)
	}

	if len(filepaths) > 1 {
		if showAll || *showLines {
			fmt.Printf("%8d ", totalLines)
		}
		if showAll || *showWords {
			fmt.Printf("%8d ", totalWords)
		}
		if showAll || *showChars {
			fmt.Printf("%8d ", totalChars)
		}
		fmt.Printf("total\n")
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
