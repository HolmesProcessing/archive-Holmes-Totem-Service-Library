import mmap
import tempfile


MEGABYTE = 2 ** 20


class TemporaryFile(object):
    """
    Easily create a temporary file.
    Automatically destroyed after proper usage.

    The temporary file is stored in memory until it either exceeds the maximum
    size or functions like fileno() are called.
    The parameter defining the maximum in memory size is in megabytes.
    The default value is 1GB.

    Usage:
        with TemporaryFile(max_memory_size=1) as file:
            file.write(b"some content")
            file.flush()
            file.seek(0)
            print(file.read())
    """

    def __init__(self, max_memory_size=1024):
        self.max_size = max_memory_size * MEGABYTE

    def __enter__(self):
        """
        Create the temporary file in memory first, when it uses too much memory
        it is automatically relocated to the filesystem.
        """
        self.file = tempfile.SpooledTemporaryFile(max_size=self.max_size)
        return self.file

    def __exit__(self, type, value, traceback):
        """
        Cleanup temporary file.
        """
        self.file.close()


class MmapFileReader (object):
    """
    Read-only file-like object trimmed for low memory footprint.
    Reading and finding does not advance the offset, however, it is position
    aware.

    Usage:
        # open
        file = MmapFileReader("/filepath")
        # find
        start = file.find(b"needle")
        # access data, still offset 0
        file[32987000:2323493493]
        # create a subfile at offset start
        subfile = file.subfile(start)
        # find a needle somewhere after the offset, relative to the offset
        position = subfile.find("second needle")
        # adjust offset in the subfile to after the previous find
        subfile.seek_relative(position+1)
    """
    __slots__ = ["file","datamap","filesize","offset"]
    def __init__ (self, filename):
        self.file     = open(filename, "rb")
        self.datamap  = mmap.mmap(self.file.fileno(), 0, access=mmap.ACCESS_READ)
        self.offset   = 0
        self.filesize = self.datamap.size()

    def __enter__ (self):
        return self

    def __exit__ (self, type, value, traceback):
        self.close()

    def close (self):
        self.datamap.close()
        del(self.datamap)
        self.file.close()
        del(self.file)
        del(self.offset)
        del(self.filesize)
        del(self)

    # provide base functionality
    def read (self, start, stop):
        remaining = self.filesize - self.offset
        # check start value
        # lowest possible value is 0
        # highest possible value is remainingsize-1
        if start is None:
            start = 0
        if start < 0:
            start = remaining + start
        start = max(0, min(start, remaining-1))
        # check stop value
        # lowest possible value is 0
        # highest possible value is remainingsize
        if stop is None:
            stop = 0
        if stop < 0:
            stop = remaining + stop
        stop = max(start, min(stop, remaining))
        # get slice, offset dependent, position unaltered after op
        self.datamap.seek(0)
        data = self.datamap[(self.offset+start):(self.offset+stop)]
        self.datamap.seek(self.offset)
        if len(data) == 1:
            return data[0]
        return data

    def seek (self, position):
        self.offset = max(0, min(position, self.filesize-1))
        self.datamap.seek(self.offset)
    def seek_relative (self, offset):
        self.seek(self.offset + offset)

    def tell (self):
        return self.offset

    def find (self, needle):
        self.datamap.seek(0)
        result = self.datamap.find(needle, self.offset)
        self.datamap.seek(self.offset)
        if result != -1:
            result -= self.offset
        return result

    def startswith (self, needle):
        return self[0:len(needle)] == needle

    def size (self):
        return len(self)

    # extended slicing
    def __getitem__ (self, key):
        if isinstance(key, slice):
            start = key.start
            stop  = key.stop
            if not start:
                start = 0
            if not stop:
                stop = len(self)
            return self.read(start, stop)
        else:
            start = key
            if not start:
                start = 0
            return self.read(start, start+1)

    def subfile (self, start):
        class MmapFileSubReader (MmapFileReader):
            __slots__ = ["file","datamap","size","offset"]
            # lightweight subtype of LargeFileReader offering adjusted offset
            def __init__ (self, file, datamap, start, size):
                self.file     = file
                self.datamap  = datamap
                self.filesize = size
                self.offset   = max(0, min(start, size))
            def close (self):
                pass  # remove close ability
            def subfile (self, start):
                pass  # remove subfile ability
        return MmapFileSubReader(self.file, self.datamap, self.offset+start, self.filesize)

    # provide standard functions
    def __len__ (self):
        return self.filesize
