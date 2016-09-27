import unittest
import tempfile
import os
from python3.tools.files import MmapFileReader


class TemporaryFileTest(unittest.TestCase):

    def setUp(self):
        file = tempfile.NamedTemporaryFile(delete=False)
        self.filename = file.name
        # write test data, then close the file:
        self.testbytesA = b"deadbeef"
        self.testbytesB = b"beefdead"
        self.data = (
            b"a"*5 +
            b"b"*5 +
            b"c"*5 +
            b"d"*5 +
            b"a"*80 +
            self.testbytesA +
            b"b"*(100-len(self.testbytesA)) +
            self.testbytesA +
            b"b"*(100-len(self.testbytesA)) +
            self.testbytesB +
            b"c"*(100-len(self.testbytesB))
        )
        file.write(self.data)
        file.close()

    def tearDown(self):
        os.remove(self.filename)

    def test_everything(self):
        with MmapFileReader(self.filename) as reader:
            # check size
            self.assertEqual(len(reader), 400)

            # at offset 0
            self.assertEqual(reader.tell(), 0)
            self.assertEqual(reader.find(self.testbytesA), 100)
            self.assertEqual(reader.find(self.testbytesB), 300)

            # at offset 100
            reader.seek(100)
            self.assertEqual(reader.tell(), 100)
            self.assertEqual(reader.find(self.testbytesA), 0)
            self.assertEqual(reader.startswith(self.testbytesA), True)
            self.assertEqual(reader.find(self.testbytesB), 200)

            # at offset 101
            reader.seek_relative(1)
            self.assertEqual(reader.tell(), 101)
            self.assertEqual(reader.find(self.testbytesA), 99)
            self.assertEqual(reader.find(self.testbytesB), 199)

            # test remaining seeking capabilities
            reader.seek_relative(-10)
            self.assertEqual(reader.tell(), 91)
            reader.seek_relative(-100)
            self.assertEqual(reader.tell(), 0)
            reader.seek(400)
            self.assertEqual(reader.tell(), 399)
            reader.seek(-1)
            self.assertEqual(reader.tell(), 0)

            # test data reading
            self.assertEqual(reader[:],       self.data)
            self.assertEqual(reader[:10],     self.data[:10])
            self.assertEqual(reader[2:20],    self.data[2:20])
            self.assertEqual(reader[-10],     self.data[-10])
            self.assertEqual(reader[-10:-5],  self.data[-10:-5])
            self.assertEqual(reader[-10:-10], self.data[-10:-10])
            self.assertEqual(reader[-10:-11], self.data[-10:-11])
            self.assertEqual(reader[2:-11],   self.data[2:-11])


if __name__ == '__main__':
    unittest.main()
