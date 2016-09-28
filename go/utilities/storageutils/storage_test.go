package storageutils

import (
	// Imports for testing
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"bytes"
	"net/http"
)

// Run some tests on the storageutils package.
var (
	storage_prot string = "http://"
	storage_addr string = "127.0.0.1"
	storage_port int    = 8017
	storage      *Storage
	fileBytes    []byte = []byte("hello world!")
	fileName     string = "testfile.txt"
	fileDate     string = time.Now().Format(time.RFC3339)
	sha256string string
)

func TestStorageutils(T *testing.T) {
	T.Log("=== Testing module storageutils ===")
	go LaunchServer(storage_addr, &storage_port)
	time.Sleep(100 * time.Millisecond)
	storage = &Storage{
		fmt.Sprintf("%s%s:%d", storage_prot, storage_addr, storage_port),
		"user-1",
	}

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

/*
Launch a mini webserver that returns some static data
(for self contained testing)
*/
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"ResponseCode":1,"Failure":""}`))
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, "", time.Now(), bytes.NewReader(fileBytes))
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetHandler(w, r)

	} else if r.Method == "PUT" {
		SubmitHandler(w, r)

	} else {
		http.NotFound(w, r)
	}
}

func LaunchServer(addr string, port *int) {
	http.HandleFunc("/samples/", RequestHandler)
	fmt.Println(http.ListenAndServe(fmt.Sprintf("%s:%d", addr, *port), nil))
}
