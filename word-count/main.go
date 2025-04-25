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

	file, err := os.Open(path)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("%v: permission denied", path)
		}
	}
	defer file.Close()
	return nil
}

func countLines(path string) (int, error) {
	if err := validateFilePath(path); err != nil {
		return 0, err
	}
	return countWithScanner(path, bufio.ScanLines)
}

func countWords(path string) (int, error) {
	if err := validateFilePath(path); err != nil {
		return 0, err
	}
	return countWithScanner(path, bufio.ScanWords)
}

func countCharacters(path string) (int, error) {
	if err := validateFilePath(path); err != nil {
		return 0, err
	}
	return countWithScanner(path, bufio.ScanRunes)
}

func countWithScanner(path string, split bufio.SplitFunc) (int, error) {
	if err := validateFilePath(path); err != nil {
		return 0, err
	}

	file, err := os.Open(path)
	if err != nil {
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

func run(args []string) error {
	flags := flag.NewFlagSet("wc", flag.ContinueOnError)
	showLines := flags.Bool("l", false, "Count lines")
	showWords := flags.Bool("w", false, "Count words")
	showChars := flags.Bool("c", false, "Count characters")
	if err := flags.Parse(args); err != nil {
		return err
	}

	remainingArgs := flags.Args()
	if len(remainingArgs) == 0 {
		return fmt.Errorf("please provide a file path")
	}
	filepath := remainingArgs[0]

	var lines, words, chars int
	var err error
	showAll := !*showLines && !*showWords && !*showChars
	if showAll || *showLines {
		lines, err = countLines(filepath)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	if showAll || *showWords {
		words, err = countWords(filepath)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	if showAll || *showChars {
		chars, err = countCharacters(filepath)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	if showAll || *showLines {
		fmt.Printf("%8d ", lines)
	}
	if showAll || *showWords {
		fmt.Printf("%8d ", words)
	}
	if showAll || *showChars {
		fmt.Printf("%8d ", chars)
	}

	fmt.Printf(" %s\n", filepath)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
