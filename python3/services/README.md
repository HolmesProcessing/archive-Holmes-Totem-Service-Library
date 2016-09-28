## Contents:
- [Prerequisites](#prerequisites)
- [JSON Configuration Parsing](#json-configuration-parsing)
- [Input Identification and Validation](#input-identification-and-validation)
- [HTTP-Router for Standard Service URL-Endpoints](#http-router-for-standard-service-url-endpoints)
- [Standardized Info-Output](#standardized-info-output)


## Prerequisites

- The [Tornado](http://www.tornadoweb.org/) web framework is installed
  ```shell-script
  sudo pip3 install tornado
  ```

- The library is on the import path.
  To add to the import path, you can e.g. do
  something like this (assuming the library is in the folder `./holmeslibrary`):
  ```python
  import os
  import sys

  dir_path = os.path.dirname(os.path.realpath(__file__))
  sys.path.append(os.path.abspath(os.path.join(dir_path, "holmeslibrary")))
  ```


## JSON Configuration Parsing

### Import
```python
from python3.services.configuration import ParseConfig
```

### Load File
Lets assume a file named `service.conf` exists, containing the following valid
JSON:

```json
{
    "Port": 8080,
    "IP": "0.0.0.0"
}
```

Then the corresponding config dictionary should look something like this with
default values:
```python
config = {
    "Port": 7777,
    "IP": "127.0.0.1",
    "IP2": "127.0.1.1", # another IP for demonstration of default values
}
```
**NOTE: The contents of `config` will be overwritten by contents specified in the JSON file.**

The library function can be used to easily get access to its contents:
```python
cfg = ServiceConfig(config, "service.conf")
print(cfg.port)
print(cfg.ip)
print(cfg.ip2)
```
The attribute access is case insensitive, so `cfg.IP2` will work too.

Since the default path is `service.conf` you could omit the file path and just
do

```python
cfg = ServiceConfig(config)
```

### Load String
The function can also be used to load the contents of a string containing valid
JSON data.

```python
jsonString = """
{
    "Port": 8080,
    "IP": "0.0.0.0"
}
"""
cfg = ServiceConfig(config, data=jsonString)
```


## Input Identification and Validation

Allows identification and validation of the following types:
- IP (v4/v6, decimal (<32 bits is automatically assumed IPv4))
- IPNet (v4/v6, no validation)
- Domain
- Email address
- File on local storage (relative to /tmp)

All errors returned are integers and can be compared to their respective values
defined in `inputtype.Errors` (see [inputtype.py#L251](https://github.com/HolmesProcessing/Holmes-Totem-Service-Library/blob/master/python3/services/inputtype.py#L251)).

Returns `inputtype.Errors.UnknownTypeError` if it cannot detect the type or
`inputtype.Errors.EmptyInputError` if the given input is an empty string.

**NOTE: For validation of Domains and Emails you must first initialize the TLD map**
(see example below for details).

Validation can return a couple different errors:
- `inputtype.ValidateIP(ip *net.IP)`:
  - `inputtype.Errors.IPisLoopbackError`
  - `inputtype.Errors.IPisUnspecifiedError`
  - `inputtype.Errors.IPisNotPublicError` (see [inputtype.py#L277](https://github.com/HolmesProcessing/Holmes-Totem-Service-Library/blob/master/python3/services/inputtype.py#L277) for list of filtered IPs)
- `inputtype.ValidateDomain(domain string)`:
  - `inputtype.Errors.InvalidDomainError`
  - `inputtype.Errors.InvalidTLDError`
- `inputtype.ValidateEmail(mail *mail.Address)`:
  - `inputtype.Errors.InvalidEmailError`
  - `inputtype.Errors.InvalidDomainError`
  - `inputtype.Errors.InvalidTLDError`
- `inputtype.ValidateFile(filepath string)`:
  - `inputtype.Errors.FileAccessDeniedError`
  - `inputtype.Errors.FileNotFoundError`

### Import
```python
from python3.services.inputtype import (
    Detect,
    InitializeTLDMap,
    ValidateIP,
    ValidateDomain,
    ValidateEmail,
    ValidateFile,
    Errors,            # Enum  Errors
    Types,             # Enum  Types
    Email,             # Class Email
)
```

### Download TLD Map
To use the validation functionality to its full extent, download a valid TLD
list like the one published by [**iana**](http://data.iana.org/TLD/tlds-alpha-by-domain.txt).
```shell-script
wget -O iana-tld-list.txt "http://data.iana.org/TLD/tlds-alpha-by-domain.txt"
```

### Detect Input Type
`Detect(_input)` accepts bytes, str or int.
```python
objs = [
    "www.google.com",
    "test@yahoo.com",
    "127.0.0.20",
]
for obj in objs:
    detected_type, parsed_object, err = Detect(obj)
    print(detected_type, parsed_object, err)
```

`parsed_object` can be one of:
- `ipaddress.IPv4Address` / `ipaddress.IPv6Address`
- `ipaddress.IPv4Network` / `ipaddress.IPv6Network`
- `str` (domain)
- `inputtype.Email`
- `str` (file, cleaned path)


### Validate Type
```python
# Initializing the TLD Map is a must, otherwise validation functions relying on
# it will raise an exception.
InitializeTLDMap("iana-tld-list.txt")

# Validate IP
ok, err = ValidateIP("127.0.0.20")
if not ok:
    if err == Errors.IPisLoopbackError:
        print("Oh no, that was the loopback")
    else:
        print(err)
else:
    print("Yay it is valid!")
```


## HTTP-Router for Standard Service URL-Endpoints
```python
from python3.services.router import Router
from python3.services.configuration import Metadata

import tornado.web

m = Metadata(
  name="test-service",
  version="1.0",
  description="some fancy description",
  copyright="you can copy as much as you like",
  license="provided without any license"
)

class AnalysisHandler(tornado.web.RequestHandler):
  def get(self):
    self.write("Hello I'm analyzing your input: {}!".format(self.get_argument("obj", strip=False)))

router = Router(metadata=m, handlers={
  "analyze": AnalysisHandler()
})
router.ListenAndServe(8080)

```


## Standardized Info-Output
```python
from python3.services.router import InfoHandler
from python3.services.configuration import Metadata

import tornado
from tornado import web, httpserver, ioloop

m = Metadata(
  name="test-service",
  version="1.0",
  description="some fancy description",
  copyright="you can copy as much as you like",
  license="provided without any license"
)

infoHandler = InfoHandler(metadata=m)

class Application(tornado.web.Application):
  def __init__(self, infoHandler):
    handlers = [
      (r"/", infoHandler),
    ]
    settings = dict(
        template_path=os.path.join(os.path.dirname(__file__), 'templates'),
        static_path=os.path.join(os.path.dirname(__file__), 'static'),
    )
    tornado.web.Application.__init__(self, handlers, **settings)
    self.engine = None

server = tornado.httpserver.HTTPServer(Application())
server.listen(8080)
tornado.ioloop.IOLoop.instance().start()
```
