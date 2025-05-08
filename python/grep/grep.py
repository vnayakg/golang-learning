import argparse
import sys
import os
from typing import Iterable


class MyGrepError(Exception):
    pass


def grep_in_file(search_string: str, filename: str) -> list[str]:
    if not os.path.exists(filename):
        raise MyGrepError(f"{filename}: open: No such file or directory")

    if os.path.isdir(filename):
        raise MyGrepError(f"{filename}: read: Is a directory")

    try:
        with open(filename, "r", encoding="utf-8") as file:
            return [line.rstrip("\n") for line in file if search_string in line]
    except PermissionError:
        raise MyGrepError(f"{filename}: Permission denied")
    except Exception as e:
        raise MyGrepError(f"{filename}: {str(e)}")


def grep_in_stdin(search_string: str, lines: Iterable[str]) -> list[str]:
    return [line.rstrip("\n") for line in lines if search_string in line]


def write_output_to_file(output_lines: str, output_file: str):
    if os.path.exists(output_file):
        raise MyGrepError(f"{output_file}: File already exists")
    try:
        with open(output_file, "w", encoding="utf-8") as f:
            for line in output_lines:
                f.write(f"{line}\n")
    except Exception as e:
        raise MyGrepError(f"{output_file}: Could not write: {str(e)}")


def parse_args():
    parser = argparse.ArgumentParser(description="A simple grep-like tool.")
    parser.add_argument("search_string", help="The string to search for")
    parser.add_argument(
        "input_file",
        nargs="?",
        help="File to search in (optional; uses stdin if not provided)",
    )
    parser.add_argument(
        "-o",
        "--output",
        type=str,
        nargs="?",
        help="Output file (optional; uses stdin if not provided)",
    )

    return parser.parse_args()


def main():
    args = parse_args()
    try:
        if args.input_file:
            matches = grep_in_file(args.search_string, args.input_file)
        else:
            matches = grep_in_stdin(args.search_string, sys.stdin)

        if args.output:
            write_output_to_file(matches, args.output)
        else:
            for line in matches:
                print(line)

    except MyGrepError as e:
        print(f"./mygrep: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
