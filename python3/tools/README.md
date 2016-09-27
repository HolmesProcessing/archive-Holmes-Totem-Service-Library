## Prerequisites

- The library is on the import path.
  To add to the import path, you can e.g. do
  something like this (assuming the library is in the folder `./holmeslibrary`):
  ```python
  import os
  import sys

  dir_path = os.path.dirname(os.path.realpath(__file__))
  sys.path.append(os.path.abspath(os.path.join(dir_path, "holmeslibrary")))
  ```

- The [Requests](http://docs.python-requests.org/en/latest/user/install/#install)
  library is installed
  ```shell-script
  sudo pip3 install requests
  ```

- The [Tornado](http://www.tornadoweb.org/) web framework is installed
  ```shell-script
  sudo pip3 install tornado
  ```

- The [python-rfc3339](https://github.com/tonyg/python-rfc3339) python package
  is installed
  ```shell-script
  pip install -e git://github.com/tonyg/python-rfc3339.git#egg=rfc3339
  ```


## Overview
- [storageutils](#storageutils)
- [MmapFileReader](#mmapfilereader)
- [TemporaryFile](#temporaryfile)


## storageutils
This toolset currently only includes a function to publish samples to a
[Holmes-Storage](https://github.com/HolmesProcessing/Holmes-Storage) instance.

### Import
```python
from python3.tools.storageutils import Storage, StorageSample
import rfc3339
```

### Usage
General usage is very simple. Create a `Storage` object, create `StorageSamples`
and have `Storage` submit these.

```python
storage = Storage(address="127.0.0.1", user_id="1")
sample = StorageSample(
    filecontents=b"Hello World!",
    source="Unknown",
    name="hello_world.txt",
    date=rfc3339.now().isoformat(),
    tags=["malware","bad-guy","simplistic"],
    comment="This file is a very good malware example ;)"
)

# submit to Holmes-Storage
storage.submitSample(sample)

# retrieve from Holmes-Storage
bytes = storage.getSample(sample.sha256())
```


## MmapFileReader
Easy to use file-like wrapper to quickly search large files by mapping them
into memory.


### Import
```python
from python3.tools.files import MmapFileReader
```


### Opening a File
```python
file = MmapFileReader("/filepath")
```


### Searching
Searching is always relative to the current offset and does not modify it.
To adjust the offset use [file.seek](#changing-position).
Returns the position or `-1` if not found.
```python
firstPosition = file.find(b"byte-sequence")
```


### Reading
Reading is always relative to the current offset and does not modify it.
To adjust the offset use [file.seek](#changing-position).
Reading out of bounds is ignored and only values within the file boundaries are
returned.
```python
data = file[0x2000:0x4000]
```


### Changing Position
Seeking an absolute value changes the offset to the specified value.
Seeking a relative value increases or decreases the offset by the given value.
Seeking values lower than 0 or greater than the file size is not possible,
instead in these cases 0 or the max offset will be set.
```python
file.seek(0x4000)
file.seek_relative(-0x1000)
```


### Creating a Subfile
It is possible to create child readers from the main MmapFileReader. This
creates a copy of the MmapFileReader that inherits all properties from the
original, except it cannot be used to create further Subfiles and it cannot be
closed (only deleted), closing is the priviledge of the main file.

Setting an offset that would result in an out of bounds offset is reset to the
respective minimum or maximum value.

To create a Subfile at the current location:
```python
file.subfile(0x0)
```

Or at current offset + 0x1000:
```python
file.subfile(0x1000)
```


## TemporaryFile
Easy to use temporary file wrapper.

### Import
```python
from python3.tools.files import TemporaryFile
```

### Usage
The constructor supports one parameter, `max_memory_size`.
This parameter measures the maximum in memory size of the file in Megabytes
before the file is written to the disk. The default value is 1 Gigabyte.

```python
with TemporaryFile(max_memory_size=1) as file:
    file.write(b"some content")
    file.flush()
    file.seek(0)
    print(file.read())
```

Upon leaving the `with` statement, the temporary file is destroyed.

If `file.fileno()` is called, the file is created on disk and starts behaving
like a regular temporary file.
