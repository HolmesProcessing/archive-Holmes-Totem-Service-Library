import json
from python3.tools.structs import StructDict
from tornado.web import HTTPError

def ParseConfig(config, path="service.conf", data=None):
    """
    Try opening the path, reading it all in and parsing it as json.
    If an error occures, throw a tornado.web.HTTPError (well defined
    behaviour by tornado for these).
    If parsing succeeds, update provided config dictionary.
    """
    if not isinstance(config, dict):
        raise ValueError("Invalid parameter supplied to ParseConfig(config), given {}, but expects a dict".format(type(config)))

    if data is None:
        try:
            with open(path, "r") as file:
                try:
                    loaded_config = json.loads(file.read())
                except Exception as e:
                    raise HTTPError(500, "Error parsing config file: {}".format(e), reason="Bad Service Configuration")
        except Exception as e:
            raise HTTPError(500, "Error opening config file: {}".format(e), reason="Bad Service Configuration")
    else:
        try:
            loaded_config = json.loads(data)
        except Exception as e:
            raise HTTPError(500, "Error parsing config input: {}".format(e), reason="Bad Service Configuration")

    __updateDict(config, loaded_config)
    return StructDict(config)

def __toLower(key):
    if isinstance(key, str):
        return key.lower()
    return key

def __updateDict(old, new):
    keymap = {}
    for key in old:
        keymap[__toLower(key)] = key

    for key in new:
        _key = __toLower(key)
        if _key in keymap:
            ofrag = old[keymap[_key]]
            nfrag = new[key]

            if isinstance(ofrag, dict):
                if not isinstance(nfrag, dict):
                    raise ValueError("Mismatching config, expected dict, got: {}".format(type(nfrag)))
                __updateDict(ofrag, nfrag)

            else:
                # other entries are replaced if their types match
                if type(ofrag) != type(nfrag):
                    raise ValueError("Mismatching config, expected {}, got: {}".format(type(ofrag),type(nfrag)))
                old[keymap[_key]] = nfrag

# Make alias
ParseConfiguration = ParseConfig


class Metadata(object):
    def __init__(self, name, version, description, copyright, license):
        self.name = name
        self.version = version
        self.description = description
        self.copyright = copyright
        self.license = license
        # TODO automatic config loading ability?
