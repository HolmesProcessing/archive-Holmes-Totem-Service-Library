## Usage

#### Import
```go
import "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/utilities/httputils"
```

#### Upload Multipart File
```go
files := []*httputils.RequestFile{
    {
        ParameterName: "file",
        FilePath:      "hello_world.txt",
        FileContent:   []byte("Hello World File!")
    },
}
request := &httputils.Request{
    URL:    "127.0.0.1:8080",
    // "POST" or "GET" - does only influence the way that parameters are stored
    Method: "POST",
    Files:  files,
}
response, err := request.Run()
```

#### Send POST Request with Parameters
```go
parameters := map[string]string{
    "param1": "value1",
    "param2": "value2",
}
request := &httputils.Request{
    URL:        "127.0.0.1:8080",
    Method:     "POST",
    Parameters: parameters
}
response, err := request.Run()
```
