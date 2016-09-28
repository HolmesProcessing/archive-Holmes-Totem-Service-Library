import unittest
import rfc3339
from python3.tools.storageutils import Storage, StorageSample

import time

import tornado.ioloop
import tornado.web


def CreateTestServer(data):
    class TestServer(tornado.web.RequestHandler):
        def get(self):
            self.set_header("Content-Type", "application/octet-stream")
            self.write(data)

        def put(self):
            self.write('{"ResponseCode":1,"Failure":""}')

    return TestServer


class StorageutilsTest(unittest.TestCase):

    def setUpClass():
        # launch webserver
        app = tornado.web.Application([(r"/samples/.*", CreateTestServer(b"hello world!"))])
        app.listen(8017)

    def tearDownClass():
        # shutdown webserver
        time.sleep(0.5)
        tornado.ioloop.IOLoop.instance().stop()

    def setUp(self):
        self.storage = Storage(
            address="http://127.0.0.1:8017",
            user_id="user-1"
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

    def test_0_submit(self):
        self.storage.submitSample(self.sample)

    def test_1_get(self):
        r = self.storage.getSample(self.sha256)
        self.assertTrue(r == self.sample.filecontents)


if __name__ == '__main__':
    unittest.main()
