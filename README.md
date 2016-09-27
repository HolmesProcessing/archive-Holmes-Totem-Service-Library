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
* Some tests require a Holmes-Storage connection (expected at 127.0.0.1:8016).
  If you want a dummy implementation of Holmes-Storage that does not store
  anything on the local hard drive, check out
  [Holmes-Storage-Testdummy](https://github.com/ms-xy/Holmes-Storage-Testdummy)

### To run the Python unit test:
```shell-script
cd holmeslibrary/python3
python3 -m unittest -v testing/*.py
```

### To run the Go unit test:
```shell-script
go get "github.com/HolmesProcessing/Holmes-Totem-Service-Library"
go test "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/utilities/httputils"
go test "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/utilities/storageutils"
go test "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/services/configuration"
go test "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/services/inputtype"
go test "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/services/router"
```

### Convenience Makefile For Testing
In order to use this Makefile, several prerequisites need to be fullfilled:
- Create a new folder
  - Put the Makefile below into it
  - Create the subfolders `src/github.com/HolmesProcessing/Holmes-Totem-Service-Library`
  - Run `git clone "github.com/HolmesProcessing/Holmes-Totem-Service-Library" "src/github.com/HolmesProcessing/Holmes-Totem-Service-Library"`

The Makefile will get, compile and execute the following programs required for testing:
- [go-daemon](https://github.com/ms-xy/go-daemon) (Required to run a Holmes-Storage-Testdummy from the Makefile)
- [Holmes-Storage-Testdummy](https://github.com/ms-xy/Holmes-Storage-Testdummy)

Further the following project dependencies are fetched automatically:
- [github.com/julienschmidt/httprouter](github.com/julienschmidt/httprouter) (Go dependency)

```makefile
gopath=$(shell pwd)
go=GOPATH=$(gopath) go
project="github.com/HolmesProcessing/Holmes-Totem-Service-Library"

default: get-dependencies run-tests

get-dependencies:
    @echo "\033[31m==-- Get Dependencies: --==\033[0m"
    $(go) get "github.com/julienschmidt/httprouter"

    [ -d "src/github.com/ms-xy/go-daemon" ] || \
        git clone "https://github.com/ms-xy/go-daemon.git" "src/github.com/ms-xy/go-daemon"
    [ -x "god" ] \
        || (cd "src/github.com/ms-xy/go-daemon" && make && cp "god" "../../../../god")

    $(go) get "github.com/ms-xy/Holmes-Storage-Testdummy"
    [ -x "Holmes-Storage-Testdummy" ] \
        || $(go) build -o "Holmes-Storage-Testdummy" "github.com/ms-xy/Holmes-Storage-Testdummy"
    [ -f "holmes-storage.conf" ] \
        || cp $(gopath)"/src/github.com/ms-xy/Holmes-Storage-Testdummy/config/storage.conf.example" "./holmes-storage.conf"
    @echo ""

run-tests:
    @echo "\033[33m> Launching Holmes-Storage Test-Server (non-persistent, in-memory)\033[0m"
    @touch holmes-storage.log
    @rm holmes-storage.log
    @./god --nohup --logfile holmes-storage.log --pidfile holmes-storage.pid -- ./Holmes-Storage-Testdummy --config=holmes-storage.conf

    @${MAKE} run-tests-helper || echo "\033[33m> Tests failed\033[0m"

    @echo "\033[33m> Stopping Holmes-Storage Test-Server\033[0m"
    @./god --stop --pidfile holmes-storage.pid
    @echo ""

run-tests-helper:
    @echo "\033[31m==-- Run Go Tests: --==\033[0m"
    @$(go) test $(project)/go/utilities/httputils
    @$(go) test $(project)/go/utilities/storageutils
    @$(go) test $(project)/go/services/configuration
    @$(go) test $(project)/go/services/inputtype
    @$(go) test $(project)/go/services/router
    @echo ""

    @echo "\033[31m==-- Run Python3 Tests: --==\033[0m"
    @cd src/$(project)/python3 && python3 -m unittest -v testing/*.py
    @echo ""
```
