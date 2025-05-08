import io
import os
import subprocess
import sys
import unittest
import tempfile
from pathlib import Path
from grep import (
    grep_in_file,
    MyGrepError,
    grep_in_stdin,
    grep_recursive,
    write_output_to_file,
)


class TestGrep(unittest.TestCase):

    def setUp(self):
        self.test_dir = tempfile.TemporaryDirectory()
        self.test_file_path = Path(self.test_dir.name) / "sample.txt"
        self.test_file_path.write_text(
            "This is the first line.\n"
            "I found the search_string in the file.\n"
            "Another line also contains the search_string.\n"
            "This line does not have anything.\n"
            "Another line also contains the Search_String.\n"
        )

    def tearDown(self):
        self.test_dir.cleanup()

    def test_no_match(self):
        matches = grep_in_file("not_present", str(self.test_file_path))

        self.assertEqual(matches, [])

    def test_multiple_matches(self):
        matches = grep_in_file("search_string", str(self.test_file_path))

        self.assertEqual(len(matches), 2)
        self.assertIn("I found the search_string in the file.", matches)
        self.assertIn("Another line also contains the search_string.", matches)

    def test_multiple_matches_case_insensitive(self):
        matches = grep_in_file("search_string", str(self.test_file_path), False)

        self.assertEqual(len(matches), 3)
        self.assertIn("I found the search_string in the file.", matches)
        self.assertIn("Another line also contains the search_string.", matches)
        self.assertIn("Another line also contains the Search_String.", matches)

    def test_file_not_exist(self):
        with self.assertRaises(MyGrepError) as context:
            grep_in_file("anything", "nonexistent.txt")

        self.assertIn("No such file or directory", str(context.exception))

    def test_is_directory(self):
        with self.assertRaises(MyGrepError) as context:
            grep_in_file("anything", self.test_dir.name)

        self.assertIn("Is a directory", str(context.exception))

    def test_permission_denied(self):
        restricted_file = Path(self.test_dir.name) / "restricted.txt"
        restricted_file.write_text("secret content\nsearch_string\n")
        restricted_file.chmod(0)

        with self.assertRaises(MyGrepError) as context:
            grep_in_file("search_string", str(restricted_file))
            self.assertIn("Permission denied", str(context.exception))
        restricted_file.chmod(0o644)

    def test_stdin_search_matches(self):
        test_input = io.StringIO("bar\nbarbazfoo\nFoobar\nfood\n")
        original_stdin = sys.stdin
        sys.stdin = test_input

        result = grep_in_stdin("foo")

        self.assertEqual(result, ["barbazfoo", "food"])
        sys.stdin = original_stdin

    def test_stdin_search_matches_case_insensitive(self):
        test_input = io.StringIO("bar\nbar\nfoo\nFoobar\nFoo\n")
        original_stdin = sys.stdin
        sys.stdin = test_input

        result = grep_in_stdin("foo", False)

        self.assertEqual(result, ["foo", "Foobar", "Foo"])
        sys.stdin = original_stdin

    def test_stdin_no_match(self):
        test_input = io.StringIO("bar\nbaz\nboo\n")
        sys.stdin = test_input

        result = grep_in_stdin("random")

        self.assertEqual(result, [])
        sys.stdin = sys.__stdin__

    def test_output_file_exists(self):
        output_lines = "bar\nbarbazfoo\nFoobar\nfood\n"

        with self.assertRaises(MyGrepError) as context:
            write_output_to_file(output_lines, self.test_file_path)
            self.assertIn("File already exists", str(context.exception))

    def test_output_file_success(self):
        lines = ["match 1", "match 2"]
        with tempfile.NamedTemporaryFile(delete=False) as tmpfile:
            tmpfile_path = tmpfile.name

        with self.assertRaises(MyGrepError) as ctx:
            write_output_to_file(lines, tmpfile_path)
        self.assertIn("File already exists", str(ctx.exception))

        os.remove(tmpfile_path)


class TestGrepInDirectory(unittest.TestCase):
    def setUp(self):
        self.temp_dir = tempfile.TemporaryDirectory()
        self.root = self.temp_dir.name

        self.file1_path = os.path.join(self.root, "file1.txt")
        with open(self.file1_path, "w") as f:
            f.write("this is a test\nno match here\nanother test line\n")

        nested_dir = os.path.join(self.root, "inner")
        os.makedirs(nested_dir)

        self.file2_path = os.path.join(nested_dir, "file2.txt")
        with open(self.file2_path, "w") as f:
            f.write("deep test line\nunrelated\n")

        self.file3_path = os.path.join(self.root, "empty.txt")
        with open(self.file3_path, "w") as f:
            f.write("nothing interesting\n")

    def tearDown(self):
        self.temp_dir.cleanup()

    def test_recursive_grep_matches(self):
        expected = [
            f"{self.file1_path}:this is a test",
            f"{self.file1_path}:another test line",
            f"{self.file2_path}:deep test line",
        ]
        result = grep_recursive("test", self.root)
        self.assertEqual(sorted(result), sorted(expected))

    def test_recursive_grep_case_sensitive(self):
        result = grep_recursive("Test", self.root)
        self.assertEqual(result, [])

    def test_recursive_grep_case_insensitive(self):
        expected = [
            f"{self.file1_path}:this is a test",
            f"{self.file1_path}:another test line",
            f"{self.file2_path}:deep test line",
        ]
        result = grep_recursive("Test", self.root, is_case_sensitive=False)
        self.assertEqual(sorted(result), sorted(expected))


if __name__ == "__main__":
    unittest.main()
