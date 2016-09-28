## Prerequisites

- [Holmes-Storage](https://github.com/HolmesProcessing/Holmes-Storage) running
  somewhere accessible (in the example on the local node at 127.0.0.1:8080)


## Usage

```go
import (
    "crypto/sha256"
    "fmt"
    "io/ioutil"

    "github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/utilities/storageutils"
)

func check (err error) {
    if err != nil {
        panic(err.Error())
    }
}

func main () {
    storage := &storageutils.Storage{
        Address: "http://127.0.0.1:8016",
        UserID: "1",
    }
    sample := &storageutils.StorageSample{
        // specify either FilePath or FileContents
        // (if FileContents exists FilePath is ignored)
        FilePath:     "testfile.txt",
        FileContents: []byte("hello world\n"),

        Source:       "Unknown",
        Name:         "testfile.txt",
        Date:         time.Now().Format(time.RFC3339),
        Tags:         []string{"malware", "nasty", "hard-to-remove"},
        Comment:      "What a dangerous file!",
    }
    err = storage.SubmitSample(sample)
    check(err)

    hSHA256 := sha256.New()
    hSHA256.Write(sample.FileContents)
    sha256string = fmt.Sprintf("%x", hSHA256.Sum(nil))

    bytes, err := storage.GetSample(sha256string)
    check(err)
    fmt.Println(string(bytes))
}
```
