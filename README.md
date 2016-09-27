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
