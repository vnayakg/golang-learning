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
			return fmt.Errorf("file does not exist: %w", err)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("file permission error: %w", err)
		}
		return fmt.Errorf("error stating file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory: %v", path)
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("file permission error: %w", err)
		}
		return fmt.Errorf("file open error: %w", err)
	}
	defer file.Close()
	return nil
}

func countLinesFromFile(filepath string) (int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, fmt.Errorf("error opening file %w", err)
	}
	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error scanning file %w", err)
	}
	return lineCount, nil
}

func countLines(path string) (int, error) {
	if err := validateFilePath(path); err != nil {
		return 0, err
	}
	return countLinesFromFile(path)
}

func run(args []string) error {
	flags := flag.NewFlagSet("wc", flag.ContinueOnError)
	lineFlag := flags.Bool("l", false, "get line count")

	if err := flags.Parse(args); err != nil {
		return err
	}

	remainingArgs := flags.Args()
	if len(remainingArgs) == 0 {
		return fmt.Errorf("please provide a file path")
	}
	filePath := remainingArgs[0]

	if *lineFlag {
		count, err := countLines(filePath)
		if err != nil {
			return fmt.Errorf("error counting lines: %w", err)
		}
		fmt.Printf("%d %s\n", count, filePath)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
