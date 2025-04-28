package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type FileResult struct {
	path                string
	lines, words, chars int
	err                 error
}

type Options struct {
	showAll, showLines, showWords, showChars bool
}

const BUF_SIZE = 1024 * 1024

func countFromReader(reader *bufio.Reader) (lines, words, chars int, err error) {
	var inWord bool

	buf := make([]byte, BUF_SIZE)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			data := buf[:n]
			chars += n
			for _, b := range data {
				if b == '\n' {
					lines++
				}
				if isSpace(b) {
					if inWord {
						inWord = false
					}
				} else {
					if !inWord {
						words++
						inWord = true
					}
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, 0, err
		}
	}
	return lines, words, chars, nil
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\v' || b == '\r'
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

func countWordsFrom(s string) int {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)
	words := 0

	for scanner.Scan() {
		words++
	}
	return words
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

func processFiles(filepaths []string, options *Options) error {
	totalLines, totalWords, totalChars := 0, 0, 0
	fileResults := make([]FileResult, len(filepaths))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for index, filepath := range filepaths {
		wg.Add(1)
		go func() error {
			if err := validateFilePath(filepath); err != nil {
				return err
			}

			file, err := os.Open(filepath)
			if err != nil {
				return fmt.Errorf("error opening file: %w", err)
			}
			defer file.Close()
			reader := bufio.NewReader(file)
			lines, words, chars, err := countFromReader(reader)

			if err != nil {
				return err
			}
			mu.Lock()
			totalLines += lines
			totalWords += words
			totalChars += chars
			mu.Unlock()

			fileResults[index] = FileResult{path: filepath, lines: lines, words: words, chars: chars}
			defer wg.Done()
			return nil
		}()
	}
	wg.Wait()

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
