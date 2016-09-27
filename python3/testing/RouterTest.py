import unittest
from python3.services.router import Router
from python3.services.configuration import Metadata

import tornado.web
import threading
import requests
import time


class TServer(threading.Thread):
    def __init__(self, metadata, analysisHandler, address):
        self.router = Router(metadata=metadata, handlers={
            "analyze": analysisHandler
        })
        self.address = address
        threading.Thread.__init__(self)
        self.daemon  = True
    def run(self):
        self.router.ListenAndServe(self.address)


class TestServiceConfig(unittest.TestCase):

    def test_all(self):
        exampleMetadata = Metadata(
            name="test-service",
            version="1.0",
            description="some fancy description",
            copyright="you can copy as much as you like",
            license="provided without any license"
        )

        class AnalysisHandler(tornado.web.RequestHandler):
            def get(self):
                self.write("Hello I'm analyzing your input: {}!".format(self.get_argument("obj", strip=False)))

        port = 7777
        address = "http://127.0.0.1:"+str(port)

        server = TServer(exampleMetadata, AnalysisHandler, port)
        server.start()
        time.sleep(0.5)

        info    = requests.get(address+"/")
        analyze = requests.get(address+"/analyze/", params={"obj":"IT'S FREAKY"})

        self.assertEqual(info.text, """
<p>test-service - 1.0</p>
<hr>
<p>some fancy description</p>
<hr>
<p>provided without any license</p>
<hr>
<p>you can copy as much as you like</p>
        """.strip())
        self.assertEqual(analyze.text, "Hello I'm analyzing your input: IT'S FREAKY!")


if __name__ == '__main__':
    unittest.main()
