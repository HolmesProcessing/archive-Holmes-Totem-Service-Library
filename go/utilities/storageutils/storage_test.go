package storageutils

import (
	// Imports for testing
	"crypto/sha256"
	"fmt"
	"testing"
	"time"
)

// Run some tests on the storageutils package.
var (
	storage_addr string   = "http://127.0.0.1:8016"
	storage      *Storage = &Storage{storage_addr, "user-1"}
	fileBytes    []byte   = []byte("hello world!")
	fileName     string   = "testfile.txt"
	fileDate     string   = time.Now().Format(time.RFC3339)
	sha256string string
)

func TestStorageutils(T *testing.T) {
	T.Log("=== Testing module storageutils ===")

	// Run test cases.
	x := true &&
		testSubmitSample(T) &&
		testGetSample(T)
	if !x {
		T.Fail()
	}
}

// Test the sample submission based on the Holmes-Storage code.
func testSubmitSample(T *testing.T) bool {
	T.Log("--- Submit Sample ---")

	hSHA256 := sha256.New()
	hSHA256.Write(fileBytes)
	sha256string = fmt.Sprintf("%x", hSHA256.Sum(nil))

	sample := &StorageSample{
		FilePath:     fileName,
		FileContents: fileBytes,
		Source:       "Unknown",
		Name:         fileName,
		Date:         fileDate,
		Tags:         []string{"malware", "nasty", "hard-to-remove"},
		Comment:      "What a dangerous file!",
	}
	err := storage.SubmitSample(sample)
	if err != nil {
		T.Log("Submission failed: " + err.Error())
		return false
	}
	return true
}

func testGetSample(T *testing.T) bool {
	T.Log("--- Get Sample ---")
	fileBytesResponse, err := storage.GetSample(sha256string)
	if err != nil {
		T.Log("Getting sample failed: " + err.Error())
		return false
	}
	if !slicesEqual(fileBytesResponse, fileBytes) {
		T.Logf("output(%s)!=expected(%s)", string(fileBytesResponse), string(fileBytes))
		return false
	}
	return true
}

func slicesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	l := len(a)
	for i := 0; i < l; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
