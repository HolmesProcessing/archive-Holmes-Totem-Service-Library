# Holmes-Totem-Service-Library

This library is built to make development of services for
[Holmes-Totem](https://github.com/HolmesProcessing/Holmes-Totem) easier.

Contained within this library is a standardized set of helpers for
[Python3](https://www.python.org/download/releases/3.0/) and
[Go](https://golang.org/).


## Available functionality for services

- Go ([examples](https://github.com/HolmesProcessing/Holmes-Totem-Service-Library/tree/master/go/services/)):
  - JSON configuration parsing
  - Input identification and validation
  - HTTP router for standard paths / Standard Info-Output
  - Submit/get samples to/from a [Holmes-Storage](https://github.com/HolmesProcessing/Holmes-Storage) instance
    ([example see here](https://github.com/HolmesProcessing/Holmes-Totem-Service-Library/tree/master/go/utilities/storageutils))

- Python3 ([examples](https://github.com/HolmesProcessing/Holmes-Totem-Service-Library/tree/master/go/services/)):
  - JSON configuration parsing
  - Input identification and validation
  - HTTP router for standard paths / Standard Info-Output
  - Submit/get samples to/from a [Holmes-Storage](https://github.com/HolmesProcessing/Holmes-Storage) instance
    ([example see here](https://github.com/HolmesProcessing/Holmes-Totem-Service-Library/tree/master/python3/tools))


## Additionally the following convenience functionality is available

- Go:
  - [httputils.Request](https://github.com/HolmesProcessing/Holmes-Totem-Service-Library/tree/master/go/utilities/httputils)

- Python3:
  - [holmeslibrary.python3.tools.files: TemporaryFile, MmapFileReader](https://github.com/HolmesProcessing/Holmes-Totem-Service-Library/tree/master/python3/tools)


## Testing

### Notes
* Some tests create temporary files on the hard drive for testing
* Some tests open a dummy webserver on 127.0.0.1:8017

### To run the Python unit test:
```shell-script
cd holmeslibrary/python3
python3 -m unittest -v testing/*.py
```

### To run the Go unit test:
```shell-script
go get "github.com/HolmesProcessing/Holmes-Totem-Service-Library"
go test "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/..."
```
