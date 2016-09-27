import unittest
import rfc3339
from python3.tools.storageutils import Storage, StorageSample


class StorageutilsTest(unittest.TestCase):

    def setUp(self):
        self.storage = Storage(
            address="http://127.0.0.1:8016",
            user_id="user-2"
        )
        self.sample = StorageSample(
            filecontents=b"hello world!",
            source="Unknown",
            name="testfile.txt",
            date=rfc3339.now().isoformat(),
            tags=["malware","really nasty file ;)"],
            comment="just some test"
        )
        self.sha256 = self.sample.sha256()

    def tearDown(self):
        pass

    def test_0_submit(self):
        self.storage.submitSample(self.sample)

    def test_1_get(self):
        r = self.storage.getSample(self.sha256)
        self.assertTrue(r == self.sample.filecontents)


if __name__ == '__main__':
    unittest.main()
