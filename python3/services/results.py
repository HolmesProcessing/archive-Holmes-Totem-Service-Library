from tornado.web import HTTPError

class ServiceResultSet (object):
    """
    Light weight result set class.
    Usage (context: tornado.web.RequestHandler):
        To create a new result set:
        > resultset = ServiceResultSet()

        Adding a subset containing some values:
        > subset = ServiceResultSet()
        > subset.add("key1","value")
        > subset.add("key2","value")
        > resultset.add("section1",subset)

        Alternatively, subsets can be skipped, in which case new dictionaries
        will be created for them (this fails if one of the keys specified is
        anything else than undefined or a dictionary):
        > resultset.add("section2","subsection1","key1","value")

        Additionally dictionaries can be added directly:
        > resultset.add({"section3":"value","section4":"value"})

        Write json response out:
        > self.write(resultset.dict())

    Output:
        {"section1":{"key1":"value","key2":"value"},"section2":{
        "subsection1":{"key1":"value"}},"section3":"value","section4":"value"}
    """
    __slots__ = ["data", "size"]

    def __init__(self):
        self.data = {}
        self.size = 0

    def add(self, *args):
        l = len(args)
        if l == 1 and isinstance(args[0], (dict, ServiceResultSet)):
            if isinstance(args[0], dict):
                self._add_dict(args[0])
            else:
                self._add_dict(args[0].dict())
        elif l > 1:
            self._add_args(self.data, args)
        else:
            raise HTTPError(
                500,
                "ServiceResultSet.add() not defined for parameters {}".format(args),
                "Service Exception")

    def dict(self):
        return self.data

    # do not call, for internal use only
    def _add_dict (self, obj):
        for (key, val) in obj.items():
            self._add_args(self.data, [key, val])

    # do not call, for internal use only
    def _add_args (self, _dict, args):
        key = args[0]
        val = args[1:]
        if len(val) > 1:
            if not (key in _dict):
                _dict[key] = {}
            if not isinstance(_dict[key], dict):
                raise HTTPError(
                    500,
                    "Key={} is not a dict".format(key),
                    "Service Exception")
            self._add_args(_dict[key], val)
        else:
            # if val is a result set itself, get its dictionary
            val = val[0]
            if isinstance(val, ServiceResultSet):
                val = val.dict()
            # check if exists
            # if exists and if not list, make list, then append
            # if not exists set
            if not (key in _dict):
                _dict[key] = val
            elif isinstance(_dict[key], list):
                _dict[key].append(val)
            else:
                tval = _dict[key]
                _dict[key] = []
                _dict[key].append(tval)
                _dict[key].append(val)
            self.size += 1
