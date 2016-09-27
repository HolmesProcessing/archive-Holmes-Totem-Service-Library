import unittest
from python3.tools.files import TemporaryFile


class TemporaryFileTest(unittest.TestCase):

    def test_1_createWriteRead(self):
        with TemporaryFile() as file:
            file.write(b"Test content of this file")
            file.flush()
            file.seek(0)
            self.assertEqual(file.read(), b"Test content of this file")


if __name__ == '__main__':
    unittest.main()
