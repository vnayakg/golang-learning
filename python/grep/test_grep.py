import unittest
import tempfile
from pathlib import Path
from grep import grep_in_file, MyGrepError


class TestMyGrep(unittest.TestCase):

    def setUp(self):
        self.test_dir = tempfile.TemporaryDirectory()
        self.test_file_path = Path(self.test_dir.name) / "sample.txt"
        self.test_file_path.write_text(
            "This is the first line.\n"
            "I found the search_string in the file.\n"
            "Another line also contains the search_string.\n"
            "This line does not have anything.\n"
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


if __name__ == "__main__":
    unittest.main()
