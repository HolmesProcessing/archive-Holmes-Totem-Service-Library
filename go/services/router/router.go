package router

import (
	//Imports for serving on a socket and handling routing of incoming request.
	"github.com/julienschmidt/httprouter"
	"net/http"

	// Imports for accessing metadata information
	"github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/services/configuration"
)

func New(m *configuration.Metadata) *Router {
	r := &Router{
		Info: CreateInfoOutputHandler(m),
	}
	return r
}

type Router struct {
	Info    func(rw http.ResponseWriter, rq *http.Request, ps httprouter.Params)
	Analyze func(rw http.ResponseWriter, rq *http.Request, ps httprouter.Params)
	Feed    func(rw http.ResponseWriter, rq *http.Request, ps httprouter.Params)
	Check   func(rw http.ResponseWriter, rq *http.Request, ps httprouter.Params)
	Results func(rw http.ResponseWriter, rq *http.Request, ps httprouter.Params)
	Status  func(rw http.ResponseWriter, rq *http.Request, ps httprouter.Params)
}

func (this *Router) ListenAndServe(addr string) error {
	router := httprouter.New()
	if this.Info != nil {
		router.GET("/", this.Info)
	}
	if this.Analyze != nil {
		router.GET("/analyze/", this.Analyze)
	}
	if this.Feed != nil {
		router.GET("/feed/", this.Feed)
	}
	if this.Check != nil {
		router.GET("/check/", this.Check)
	}
	if this.Results != nil {
		router.GET("/results/", this.Results)
	}
	if this.Status != nil {
		router.GET("/status/", this.Status)
	}
	return http.ListenAndServe(addr, router)
}
