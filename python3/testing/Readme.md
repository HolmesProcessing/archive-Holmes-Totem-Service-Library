# Test Suite for the python3 version of the Holmes-Totem-Service-Library

## The following classes have test cases defined:

* services.configuration.ServiceConfig
* services.results.ServiceResultSet
* tools.files.TemporaryFile
* tools.files.MmapFileReader
* tools.storageutils.Storage

## Notes
* The MmapFileReader test writes some KB of test data to the
  filesystem in form of a named temporary file, which is removed again after
  execution
* In order to test the Holmes-Storage dependent classes, a Holmes-Storage server
  needs to be launched. (Expected at 127.0.0.1:8016)

## To run the test cases:

`$> cd holmeslibrary/python3`
`$> python3 -m unittest -v testing/*.py`
