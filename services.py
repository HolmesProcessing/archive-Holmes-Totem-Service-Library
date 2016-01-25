import os
import sys

# correctly import renamed modules (Py2 vs Py3)
if sys.version_info >= (3,):
    import configparser
else:
    import ConfigParser
    configparser = ConfigParser


class ServiceMeta():
    """
    Class for storing metadata for a service.
    Metadata are read from an INI style configuration file.
    Example INI:
        [service]
        ServiceName         = HelloWorld
        ServiceVersion      = 1.0
        ServiceDescription  = ./DESCRIPTION
        ServiceConfig       = ./service.conf
        ServiceCopyright    = ./COPYRIGHT
        ServiceLicense      = ./LICENSE

        [object]
        ObjectCategory      = None
        ObjectType          = None
    """
    
    needed_meta_data = [
        "ServiceName",
        "ServiceVersion",
        "ServiceDescription",
        "ServiceConfig",
        "ServiceCopyright",
        "ServiceLicense",
        
        "ObjectCategory",
        "ObjectType",
    ]

    def __init__(self, cfg="./META"):
        parser = configparser.ConfigParser()
        # avoid case insensitivity for keys
        parser.optionxform = str
        self.data = {}
        parser.read(cfg)
        # get values from any section
        # TODO: should there be treatment for specific sections?
        # - [metadata]
        # - [config]
        # ? (metadata and config could be in the same file this way)
        for section in parser.sections():
            for (key, value) in parser.items(section):
                path = False
                if value.startswith("./") or value.startswith("/"):
                    path = value
                if path and os.path.isfile(path):
                    with open(path) as file:
                        value = file.read()
                self.data[key] = value
        
        for needed in ServiceMeta.needed_meta_data:
            if self.data.get(needed) is None:
                print("%s is not configured in %s!" % (needed, cfg))
    
    def __getattr__ (self, key):
        data = self.data.get(key)
        if not data:
            data = ""
        return data
    
    def __iter__ (self):
        for key in self.data:
            yield (key, self.data[key])



class ServiceRequestError (Exception):
    """
    Basic exception class.
    Usage (context: tornado.web.RequestHandler):
       self.set_status(e.status)
       self.write(e)
    """
    __slots__ = ["status", "error"]
    def __init__ (self, status, error):
        self.status = status
        self.error  = error
    def __str__ (self):
        return str(self.status) + ": " + str(self.error)
    def __repr__ (self):
        return repr(str(self))
    def __iter__ (self):
        yield ("status", self.status)
        yield ("error", self.error)
    def __getitem__ (self, key):
        return getattr(self,key)



class ServiceResultSet (object):
    """
    Light weight result set class.
    Usage (context: tornado.web.RequestHandler):
        resultset = ResultSet()
        subset = Resultset
        subset.add("key1","value")
        subset.add("key2","value")
        resultset.add("key3",subset)
        self.write(resultset)
    Output:
        {"key3":{"key1":"value","key2":"value"}}
    """
    __slots__ = ["data", "size"]
    def __init__(self):
        self.data = {}
        self.size = 0
    def add(self, key, value):
        if key in self.data:
            if isinstance(self.data[key], list):
                self.data[key].append(value)
            else:
                cpy = self.data[key]
                self.data[key] = []
                self.data[key].append(cpy)
                self.data[key].append(value)
        else:
            self.data[key] = value
        self.size += 1


