package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type FileResult struct {
	path                string
	lines, words, chars int
	err                 error
}

type Options struct {
	showAll, showLines, showWords, showChars bool
}

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
	return countItems(path, bufio.ScanLines)
}

func countWords(path string) (int, error) {
	return countItems(path, bufio.ScanWords)
}

func countCharacters(path string) (int, error) {
	return countItems(path, bufio.ScanRunes)
}

func countWordsFrom(s string) int {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)
	words := 0

	for scanner.Scan() {
		words++
	}
	return words
}

func countItems(path string, split bufio.SplitFunc) (int, error) {
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

func parseFlags(args []string) (*Options, []string, error) {
	flags := flag.NewFlagSet("wc", flag.ContinueOnError)
	showLines := flags.Bool("l", false, "Count lines")
	showWords := flags.Bool("w", false, "Count words")
	showChars := flags.Bool("c", false, "Count characters")
	if err := flags.Parse(args); err != nil {
		return nil, nil, err
	}
	showAll := !*showLines && !*showWords && !*showChars

	options := Options{showAll: showAll,
		showLines: *showLines,
		showWords: *showWords,
		showChars: *showChars}
	return &options, flags.Args(), nil
}

func printResult(lines, words, chars int, path string, options *Options) {
	if options.showAll || options.showLines {
		fmt.Printf("%8d ", lines)
	}
	if options.showAll || options.showWords {
		fmt.Printf("%8d ", words)
	}
	if options.showAll || options.showChars {
		fmt.Printf("%8d ", chars)
	}
	fmt.Printf("%s\n", path)
}

func processSTDIN(options *Options) error {
	scanner := bufio.NewScanner(os.Stdin)
	lines, words, chars := 0, 0, 0

	for scanner.Scan() {
		line := scanner.Text()
		lines++
		words += countWordsFrom(line)
		chars += len(line) + 1
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stdin: %w", err)
	}
	printResult(lines, words, chars, "\n", options)
	return nil
}

func countFromFile(filepath string, options *Options) (*FileResult, error) {
	fileResult := &FileResult{path: filepath}
	lines, words, chars := 0, 0, 0
	var err error

	if options.showAll || options.showLines {
		lines, err = countLines(filepath)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		fileResult.lines = lines
	}
	if options.showAll || options.showWords {
		words, err = countWords(filepath)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		fileResult.words = words
	}
	if options.showAll || options.showChars {
		chars, err = countCharacters(filepath)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		fileResult.chars = chars
	}
	return fileResult, nil
}

func processFiles(filepaths []string, options *Options) error {
	totalLines, totalWords, totalChars := 0, 0, 0
	fileResults := make([]FileResult, 0, len(filepaths))

	for _, filepath := range filepaths {
		if err := validateFilePath(filepath); err != nil {
			return err
		}

		result, err := countFromFile(filepath, options)
		if err != nil {
			return err
		}
		totalLines += result.lines
		totalWords += result.words
		totalChars += result.chars

		fileResults = append(fileResults, *result)
	}

	for _, result := range fileResults {
		if result.err != nil {
			fmt.Fprintln(os.Stderr, result.err)
			continue
		}
		printResult(result.lines, result.words, result.chars, result.path, options)
	}

	if len(filepaths) > 1 {
		printResult(totalLines, totalWords, totalChars, "total", options)
	}
	return nil
}

func run(args []string) error {
	options, filepaths, err := parseFlags(args)
	if err != nil {
		return fmt.Errorf("error parsing arguments: %w", err)
	}

	if len(filepaths) == 0 {
		processSTDIN(options)
		return nil
	}
	processFiles(filepaths, options)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
