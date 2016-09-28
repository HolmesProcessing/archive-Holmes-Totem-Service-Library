import unittest
from python3.services.results import ServiceResultSet
from tornado.web import HTTPError


class TestServiceResultSet(unittest.TestCase):

    def test_1_addSubset(self):
        rset = ServiceResultSet()
        subset = ServiceResultSet()
        subset.add("key1", "value1")
        subset.add("key2", "value2")
        rset.add("section1", subset)
        rdict = rset.dict()
        self.assertEqual(rdict["section1"]["key1"], "value1")
        self.assertEqual(rdict["section1"]["key2"], "value2")

    def test_2_addMultidimensional(self):
        rset = ServiceResultSet()
        rset.add("section1","subsection1","key1","value1")
        rdict = rset.dict()
        self.assertEqual(rdict["section1"]["subsection1"]["key1"], "value1")
        with self.assertRaises(HTTPError):
            rset.add("section1","subsection1","key1","subkey1","value2")

    def test_3_addDictionary(self):
        rset = ServiceResultSet()
        rset.add({"section1": {"key1": "value1", "key2": "value2"}})
        rdict = rset.dict()
        self.assertEqual(rdict["section1"]["key1"], "value1")
        self.assertEqual(rdict["section1"]["key2"], "value2")


if __name__ == '__main__':
    unittest.main()
