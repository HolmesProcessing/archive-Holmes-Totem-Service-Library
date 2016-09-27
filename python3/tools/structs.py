
class StructDict (dict):
    """
    Extended dict class that supports access to keys via case insensitive
    attributes.

    Example:
    foo = StructDict()
    foo.x = {"hello":"world"}
    print(foo.x.hello)
    """
    def __init__ (self, data={}, keyMod=None):
        newdata = {}
        object.__setattr__(self,"__keymap", {})
        keymap = object.__getattribute__(self, "__keymap")
        for key in data:
            _key = key
            if keyMod != None:
                _key = keyMod(key)
            _key_lower = _key
            if isinstance(_key, (str, bytes)):
                _key_lower = _key.lower()
            keymap[_key_lower] = _key
            newdata[_key] = Struct(data[key])
        dict.__init__(self, newdata)

    def __getattr__ (self, key):
        if isinstance(key, (str, bytes)):
            key = key.lower()
        keymap = object.__getattribute__(self, "__keymap")
        if key in keymap:
            return self.get(keymap[key])
        return None

    def __setattr__ (self, key, value):
        _key = key
        if isinstance(key, (str, bytes)):
            _key = key.lower()
        keymap = object.__getattribute__(self, "__keymap")
        keymap[_key] = key
        self[key] = Struct(value)


def Struct (obj, keyMod=None):
    """
    Determines the type of the object and returns an appropriate wrapper if any.
    (StructDict, StructList, or just the plain input object)
    """
    if isinstance(obj, (str, int, float)):
        return obj
    if isinstance(obj, (object, dict, list, tuple)):
        return StructDict(obj, keyMod)
    else:
        return obj
