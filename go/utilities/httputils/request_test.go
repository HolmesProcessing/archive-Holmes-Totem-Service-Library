package httputils

import (
	// Imports for testing
	"io/ioutil"
	// "strings"
	"testing"
)

// Run some tests on the httputils package.
func TestHttputils(T *testing.T) {
	T.Log("=== Testing module httputils ===")
	T.Log("--- File upload request generation ---")

	parameters := map[string]string{
		"key1": "param1",
		"key2": "param2",
	}

	name := "somefilename"
	content := []byte("somefilecontents")

	request := createFileUploadRequest("127.0.0.1:8080", parameters, "sample", name, content)
	builtRequest, err := request.Build()
	if err != nil {
		T.Log(err)
		T.Fail()
	}
	builtRequest.ParseMultipartForm(0x10000)
	file, _ := builtRequest.MultipartForm.File["sample"][0].Open()
	bytes, _ := ioutil.ReadAll(file)
	output := string(bytes)
	expected := "somefilecontents"

	if output != expected {
		T.Logf("output(%v)!=%v", output, expected)
		T.Fail()
	}

	// TODO more tests (some cookie / other url schemas (containing ? e.g.))
	// A more thorough test is the storage-submit-sample_test in
	// ...holmeslib.../go/utilities/storageutils/

}
