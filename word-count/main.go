package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type FileResult struct {
	path                string
	lines, words, chars int
	err                 error
}

type Options struct {
	showLines, showWords, showChars bool
}

const BUF_SIZE = 1024 * 1024

func countFile(reader io.Reader) (lines, words, chars int, err error) {
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
					inWord = false
				} else if !inWord {
					words++
					inWord = true
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				return lines, words, chars, nil
			}
			return 0, 0, 0, err
		}
	}
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

func processFiles(paths []string) ([]*FileResult, error) {
	fileResults := make([]*FileResult, len(paths))
	var wg sync.WaitGroup

	for index, path := range paths {
		wg.Add(1)
		go func() error {
			if err := validateFilePath(path); err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("error opening file: %w", err)
			}
			defer file.Close()
			reader := bufio.NewReader(file)
			lines, words, chars, err := countFile(reader)

			if err != nil {
				return err
			}

			fileResults[index] = &FileResult{path: path, lines: lines, words: words, chars: chars}
			defer wg.Done()
			return nil
		}()
	}
	wg.Wait()

	return fileResults, nil
}

func processStdin() (*FileResult, error) {
	lines, words, chars, err := countFile(os.Stdin)

	if err != nil {
		return nil, fmt.Errorf("error reading stdin: %w", err)
	}
	return &FileResult{path: "", lines: lines, words: words, chars: chars}, nil
}

func printResult(lines, words, chars int, path string, options *Options) {
	if options.showLines {
		fmt.Printf("%8d ", lines)
	}
	if options.showWords {
		fmt.Printf("%8d ", words)
	}
	if options.showChars {
		fmt.Printf("%8d ", chars)
	}
	fmt.Printf("%s\n", path)
}

func calculateTotal(results []*FileResult) *FileResult {
	total := &FileResult{path: "total"}
	for _, result := range results {
		if result != nil {
			total.lines += result.lines
			total.words += result.words
			total.chars += result.chars
		}
	}
	return total
}

func parseFlags(args []string) (*Options, []string, error) {
	options := Options{}
	flags := flag.NewFlagSet("wc", flag.ContinueOnError)
	flags.BoolVar(&options.showLines, "l", false, "Count lines")
	flags.BoolVar(&options.showWords, "w", false, "Count words")
	flags.BoolVar(&options.showChars, "c", false, "Count characters")

	if err := flags.Parse(args); err != nil {
		return nil, nil, err
	}

	if !options.showLines && !options.showWords && !options.showChars {
		options.showLines = true
		options.showWords = true
		options.showChars = true
	}
	return &options, flags.Args(), nil
}

func run(args []string) error {
	options, paths, err := parseFlags(args)
	if err != nil {
		return fmt.Errorf("error parsing arguments: %w", err)
	}

	if len(paths) == 0 {
		result, err := processStdin()
		if err != nil {
			return fmt.Errorf("error processing stdin: %w", err)
		}
		printResult(result.lines, result.words, result.chars, "\n", options)
		return nil
	}

	fileResults, err := processFiles(paths)
	if err != nil {
		return fmt.Errorf("error processing file: %w", err)
	}

	for _, result := range fileResults {
		if result.err != nil {
			fmt.Fprintln(os.Stderr, result.err)
			continue
		}
		printResult(result.lines, result.words, result.chars, result.path, options)
	}

	if len(paths) > 1 {
		result := calculateTotal(fileResults)
		printResult(result.lines, result.words, result.chars, result.path, options)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
