package router

import (
	//Imports for serving on a socket and handling routing of incoming request.
	"github.com/julienschmidt/httprouter"
	"net/http"

	// Imports for accessing metadata information
	"github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/services/configuration"

	// Imports for serving requests
	"fmt"
)

// Create a function that outputs meta data for the information url.
// Output format:
//
//     Name - Version
//     ----------
//     Description
//     ----------
//     License
//
// Use with julienschmidt's httprouter:
//
//     m := configuration.Metadata{Name: "test", Version: "0.1"}
//     router.GET("/url/", CreateInfoOutputFunc(m))
//
func CreateInfoOutputHandler(m *configuration.Metadata) func(rw http.ResponseWriter, rq *http.Request, ps httprouter.Params) {
	return func(rw http.ResponseWriter, rq *http.Request, ps httprouter.Params) {

		fmt.Fprintf(rw, "<p>%s - %s</p>", m.Name, m.Version)

		if m.Description != "" {
			fmt.Fprintf(rw, "<hr><p>%s</p>", m.Description)
		}

		if m.License != "" {
			fmt.Fprintf(rw, "<hr><p>%s</p>", m.License)
		}

	}
}
