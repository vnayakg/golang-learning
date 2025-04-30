# word-count

Implementation of the Unix `wc` command written in Go.

## Features

- Count lines, words, and characters in text files
- Process multiple files simultaneously using worker pools
- Support for reading from stdin

## Installation
### Prerequisites

- Go 1.24.2

### Building from source
Clone this repository

```bash
cd word-count
go build
```

### Running tests
```bash
cd word-count
go test -cover
```

## Usage

### Basic usage

```bash
# Process a single file
./word-count filename.txt
#    9      9     20 filename.txt

# Process multiple files
./word-count file1.txt file2.txt
#    5      7     35 file1.txt
#   12     24    120 file2.txt
#   17     31    155 total

# Read from stdin
cat filename.txt | ./word-count
#    9      9     20
```

### Command-line flags

By default, `word-count` displays line, word, and character counts. Use these flags to customize output:

| Flag | Description |
|------|-------------|
| `-l` | Display only line count |
| `-w` | Display only word count |
| `-c` | Display only character count |

Examples:

```bash
# Count only lines
./word-count -l filename.txt
#    9 filename.txt

# Count only words and characters
./word-count -w -c filename.txt
#      9     20 filename.txt
```

## Implementation Details

- Uses a worker pool pattern with buffered channels for concurrent file processing
- Buffer size (1MB) for efficient file reading
- Processes multiple files in parallel using available CPU cores
- Handles errors gracefully with informative messages

## Performance

Performance tests conducted on 2 GB files for 50 iterations without any flags:

| Tool | Average Time |
|------|--------------|
| `wc` (standard Unix utility) | 2.876 seconds |
| `word-count` (this program) | 2.612 seconds |

Time measurements were calculated using the Unix `time` command.

## Future Optimizations

- Implement `lseek` system call for character counting (similar to original `wc`)
- Individual flags performance is not up to the original wc, work on individual flags processing

