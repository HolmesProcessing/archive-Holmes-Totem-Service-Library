##

- [JSON Configuration Parsing](#json-configuration-parsing)
- [Input Identification and Validation](#input-identification-and-validation)
- [HTTP-Router for Standard Service URL-Endpoints](#http-router-for-standard-service-url-endpoints)
- [Standardized Info-Output](#standardized-info-output)


## JSON Configuration Parsing

```go
import (
    "fmt"
    "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/services/configuration"
)

func main() {
    filepath := "config.json"

    // The configuration struct is filled with the values found in "config.json"
    var config struct{
        BindAddress string
        Section1    struct{
            Somekey string
            Subsection struct{
                Somekey string
            }
        }
    }

    // Note: Parse panics if it encounters an error.
    configuration.Parse(&config, filepath)

    // Access values just like on any struct
    fmt.Println(config.BindAddress)
}
```


## Input Identification and Validation

Allows identification and validation of the following types:
- IP (v4/v6, no decimal)
- IPNet (v4/v6, no validation)
- Domain
- Email address
- File on local storage (relative to /tmp)

All errors returned are integers and can be compared to their respective values
defined in the package `inputtype` (see [variables.go#L93](https://github.com/HolmesProcessing/Holmes-Totem-Service-Library/blob/master/go/services/inputtype/variables.go#L93)).

Returns `inputtype.UnknownTypeError` if it cannot detect the type or
`inputtype.EmptyInputError` if the given input is an empty string.
**For validation of Domains and Emails you must first initialize the TLD map**
(see example below for details).

Validation can return a couple different errors:
- `inputtype.ValidateIP(ip *net.IP)`:
  - `inputtype.IPisLoopbackError`
  - `inputtype.IPisUnspecifiedError`
  - `inputtype.IPisNonPublicError` (see [variables.go#L139](https://github.com/HolmesProcessing/Holmes-Totem-Service-Library/blob/master/go/services/inputtype/variables.go#L139) for list of filtered IPs)
- `inputtype.ValidateDomain(domain string)`:
  - `inputtype.InvalidDomainError`
  - `inputtype.InvalidTLDError`
- `inputtype.ValidateEmail(mail *mail.Address)`:
  - `inputtype.InvalidDomainError`
  - `inputtype.InvalidTLDError`
- `inputtype.ValidateFile(filepath string)`:
  - `inputtype.FileAccessDeniedError`
  - `inputtype.FileNotFoundError`


```go
import (
    "fmt"
    "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/services/inputtype"
    "io/ioutil"
    "net/http"
)

func check(err error) {
    if err != nil {
        panic(err.Error())
    }
}

func main() {
    // Note: The first two steps (downloading and initializing the tld list) are
    // only required for validating domains and email addresses.

    // Download tld list from iana.org.
    // (Don't do this in your service upon start, but rather download it in the
    // service container build, e.g. using wget in the Dockerfile)
    client := &http.Client{}
    r, _ := client.GET("http://data.iana.org/TLD/tlds-alpha-by-domain.txt")
    bytes, _ := ioutil.ReadAll(r.Body)
    ioutil.WriteFile("iana-tld-list.txt", bytes, 0644)

    // Initialize TLD map in the inputtype package.
    InitializeTLDMap("iana-tld-list.txt")

    // Hand in some test input for detection.
    input := "www.subdomain.my-domain.net"
    detectedType, domain, err := inputtype.Detect(input)
    check(err)
    fmt.Println(detectedType)
    fmt.Println(domain)

    // Validate the domain.
    _, err := inputtype.ValidateDomain()
    check(err)
}
```


## HTTP-Router for Standard Service URL-Endpoints

```go
import (
    "github.com/julienschmidt/httprouter"
    "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/services/configuration"
    "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/services/router"
    "net/http"
)

func check(err error) {
    if err != nil {
        panic(err.Error())
    }
}

func main() {
    metadata := &configuration.Metadata{
        Name:        "service-name",
        Version:     "0.1",
        Description: "service-description",
        Copyright:   "service-copyright",
        License:     "service-license",
    }
    router := router.New(metadata)
    router.Analyze = handlerAnalyze
    // Further endpoints available:
    // - router.Feed
    // - router.Check
    // - router.Results
    // - router.Status
    err := router.ListenAndServe("127.0.0.1:8080")
    check(err)
}

func handlerAnalyze(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    w.Write([]byte("Some test output for /analyze/"))
}
```


## Standardized Info-Output

This is automatically created when using the _router.New(*configuration.Metadata)_
method to create a default router.

However, it might be preferable to create your own router. In this case it is
possible to use _router.CreateInfoOutputHandler(*configuration.Metadata)_ to create
a handler for use with _httprouter_.

```go
import (
    "github.com/julienschmidt/httprouter"
    "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/services/configuration"
    "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/services/router"
    "net/http"
)

func check(err error) {
    if err != nil {
        panic(err.Error())
    }
}

func main() {
    metadata := &configuration.Metadata{
        Name:        "service-name",
        Version:     "0.1",
        Description: "service-description",
        Copyright:   "service-copyright",
        License:     "service-license",
    }
    infoHandler := router.CreateInfoOutputHandler(metadata)

    router := httprouter.New()
    router.GET("/", infoHandler)

    err := http.ListenAndServe("127.0.0.1:8080", router)
    check(err)
}
```
