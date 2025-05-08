import argparse
import sys
import os
from typing import Iterable, TextIO


class MyGrepError(Exception):
    pass


def grep(search_string: str, file: TextIO, is_case_sensitive=True) -> list[str]:
    if is_case_sensitive:
        return [line.rstrip("\n") for line in file if search_string in line]
    return [line.rstrip("\n") for line in file if search_string.lower() in line.lower()]


def grep_in_file(
    search_string: str, filename: str, is_case_sensitive=True
) -> list[str]:
    if not os.path.exists(filename):
        raise MyGrepError(f"{filename}: open: No such file or directory")

    if os.path.isdir(filename):
        raise MyGrepError(f"{filename}: read: Is a directory")

    try:
        with open(filename, "r", encoding="utf-8") as file:
            return grep(search_string, file, is_case_sensitive)
    except PermissionError:
        raise MyGrepError(f"{filename}: Permission denied")
    except Exception as e:
        raise MyGrepError(f"{filename}: {str(e)}")


def grep_in_stdin(search_string: str, is_case_sensitive=True) -> list[str]:
    return grep(search_string, sys.stdin, is_case_sensitive)


def grep_recursive(
    search_string: str, root_dir: str, is_case_sensitive=True
) -> list[str]:
    matches = []
    for dirpath, _, filenames in os.walk(root_dir):
        for filename in filenames:
            filepath = os.path.join(dirpath, filename)
            try:
                file_matches = grep_in_file(search_string, filepath, is_case_sensitive)
                if file_matches:
                    matches.extend([f"{filepath}:{line}" for line in file_matches])
            except MyGrepError:
                continue
    return matches


def write_output_to_file(output_lines: list[str], output_file: str):
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
    parser.add_argument(
        "-i",
        "--ignorecase",
        nargs="?",
        help="Perform case insensitive matching. By default, it is case sensitive",
    )
    parser.add_argument(
        "-r",
        "--recursive",
        nargs="?",
        help="Recursively search subdirectories listed",
    )

    return parser.parse_args()


def main():
    args = parse_args()
    try:
        if args.recursive:
            if not args.input_file or not os.path.isdir(args.input_file):
                raise MyGrepError("Recursive flag requires a directory path as input")
            matches = grep_recursive(
                args.search_string, args.input_file, args.ignorecase
            )
        elif args.input_file:
            matches = grep_in_file(args.search_string, args.input_file, args.ignorecase)
        else:
            matches = grep_in_stdin(args.search_string, args.ignorecase)

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
