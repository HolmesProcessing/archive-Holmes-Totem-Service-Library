import requests
import hashlib
from tornado.web import HTTPError

class StorageSample (object):
    """
    Holmes-Storage related utility class matching the Holmes-Storage sample
    structure.
    """
    def __init__ (self, filepath="", filecontents=b"", source="", name="", date="", tags=[], comment=""):
        """
        Parameters:
            filepath        - String:   Path to the file (Supply either filepath or contents)
            filecontents    - Bytes:    Contents of the file (Supply either filepath or contents)
            source          - String:   Source of the file
            name            - String:   Name of the file
            date            - String:   Date of the submission
            tags            - []String: Array of tags that are to be associated with the submission
            comment         - String:   Comment to be associated with the submission
        """
        # Local only:
        self.filepath     = filepath
        self.filecontents = filecontents

        # Contents for submission:
        self.source  = source
        self.name    = name
        self.date    = date
        self.tags    = tags
        self.comment = comment

    def getContent(self):
        if not self.filecontents:
            with open(self.filepath, "rb") as file:
                self.filecontents = file.read()
        return self.filecontents

    def sha256(self):
        return self.getHash()
    def getHash(self):
        return hashlib.sha256(self.filecontents).hexdigest()


class Storage (object):
    """
    Holmes-Storage utility wrapper class.
    """
    def __init__ (self, address, user_id):
        """
        Parameters:
            address - IP:PORT
            user_id - user id for storing data in storage
        """
        self.address = address
        self.user_id = user_id

    def submitSample (self, sample):
        """
        Parameters:
            sample - StorageSample instance
        Raises
            HTTPError - if response is malformed or if the request did fail.
        """
        url = self.address + "/samples/"
        files = {
            "sample": sample.getContent()
        }
        params = {
            "user_id": self.user_id,
            "source":  sample.source,
            "name":    sample.name,
            "date":    sample.date,
            "tags":    sample.tags,
            "comment": sample.comment
        }
        r = requests.request("PUT", url, files=files, params=params)
        try:
            r = r.json()
        except Exception as e:
            raise HTTPError(500, "Error parsing response: {}".format(e), reason="Malformed Response")

        if not "ResponseCode" in r:
            raise HTTPError(500, "Missing field 'ResponseCode': {}".format(r), reason="Malformed Response")

        if r["ResponseCode"] != 1:
            if not "Failure" in r:
                raise HTTPError(500, "Missing field 'Failure': {}".format(r), reason="Malformed Response")
            raise HTTPError(500, "Failure: {}".format(r["Failure"]), reason="Submit Failure")

    def getSample (self, sha256):
        """
        Parameters:
            sha256 - the sha256 hash of the requested samples file contents as
                     a hex string.
        Raises
            HTTPError - if response is malformed or if the request did fail.
        Returns:
            Sample contents (bytes)
        """
        url = self.address + "/samples/" + sha256
        r = requests.request("GET", url)
        if r.status_code != 200:
            raise HTTPError(500, r.content.encode("utf-8"))
        if (not r.headers) or r.headers["content-type"] != "application/octet-stream":
            try:
                r = r.json()
            except Exception as e:
                raise HTTPError(500, "Error parsing response: {}".format(e), reason="Malformed Response")

            if not "Failure" in r:
                raise HTTPError(500, "Missing field 'Failure': {}".format(r), reason="Malformed Response")

            raise HTTPError(500, "Failure: {}".format(r["Failure"]), reason="Get Failure")

        return r.content
