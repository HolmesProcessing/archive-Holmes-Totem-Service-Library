package router

import (
	// Imports for testing
	"github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/services/configuration"
	"github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/utilities/httputils"
	"io/ioutil"
	"testing"
	"time"
)

// Run some tests on the router implementation.
// Currently only testing the generation of the info-output function.
func TestRouter(T *testing.T) {
	T.Log("=== Testing module router ===")

	// Test metadata
	m := &configuration.Metadata{
		Name:        "test-service",
		Version:     "0.1",
		Description: "some-description",
		Copyright:   "some-copyright",
		License:     "some-license",
	}
	router := New(m)

	// Launch server at a test address.
	addr := "127.0.0.1:8080"
	var err error = nil
	go func() {
		err = router.ListenAndServe(addr)
	}()
	// Quarter second should be more than enough to catch any socket errors upon
	// creating the server.
	time.Sleep(250 * time.Millisecond)
	if err != nil {
		T.Log("Unexpected test server error: " + err.Error())
		T.Fail()
		return
	}

	// Create a request instance and launch a request.
	request := &httputils.Request{
		URL:         "http://" + addr,
		Method:      "GET",
		Parameters:  map[string]string{},
		Cookies:     nil,
		Files:       nil,
		Body:        nil,
		ContentType: "",
	}
	response, err := request.Run()
	if err != nil {
		T.Log("Unexpected request error: " + err.Error())
		T.Fail()
		return
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		T.Log("Unexpected I/O-Error: " + err.Error())
		T.Fail()
		return
	}

	out := string(bytes)
	expected := "<p>test-service - 0.1</p><hr><p>some-description</p><hr><p>some-license</p>"

	if out != expected {
		T.Logf("out(%s) != expected(%s)", out, expected)
		T.Fail()
	}
}
