// Package server is an https server
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Server is ...
type Server struct {
	http.Handler
}

// Response is ...
type Response struct {
}

// New initializes a new server.
func New(port string) *Server {
	srv := new(Server)

	router := httprouter.New()

	router.GET("/", srv.Index)
	// router.GET("/jobs")
	// router.GET("/jobs/:id")
	// router.POST("/jobs")
	// router.DELETE("/jobs/:id")
	// router.GET("/jobs/:id/log")

	srv.Handler = router
	return srv
}

// ***TEMP***
//
// Index handles requests to root directory
func (s *Server) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Printf("Hello from root")
	msg := map[string]string{"msg": "Hello", "name": "Brady"}

	res, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(res)

}

// get list handler /jobs
func (s *Server) listJobs() {
}

// get 1 job handler /jobs/:id
func (s *Server) jobStatus() {}

// post 1 job handler /jobs
func (s *Server) startJob() {}

// delete 1 job handler /jobs/:id
func (s *Server) stopJob() {}

// get 1 job log handler /jobs/:id/log
func (s *Server) logJob() {}
