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


def main():
    try:
        if len(sys.argv) == 3:
            search_string = sys.argv[1]
            filename = sys.argv[2]
            matches = grep_in_file(search_string, filename)
            for line in matches:
                print(line)
        elif len(sys.argv) == 2:
            search_string = sys.argv[1]
            matches = grep_in_stdin(search_string, sys.stdin)
            for line in matches:
                print(line)
        else:
            print("Usage: ./mygrep <search_string> <filename>", file=sys.stderr)
            print("Usage: ./mygrep <search_string>", file=sys.stderr)
            sys.exit(1)

    except MyGrepError as e:
        print(f"./mygrep: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
