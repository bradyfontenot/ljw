// Package server is an https server
package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	port     = ":8080"
	certFile = "ssl/server.crt"
	keyFile  = "ssl/server.key"
	caFile   = "ssl/ca.crt"
)

// Server is ...
type Server struct {
	*http.Server
}

// Response is ...
type Response struct {
}

// New initializes a new server.
func New() *Server {

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// load certificate and private key files
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	s := new(Server)

	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}

	// build router and routes
	r := httprouter.New()

	r.GET("/", Index)
	r.GET("/api/jobs", s.getJobs)
	// r.GET("/api/jobs/:id")
	// r.POST("/api/jobs")
	// r.DELETE("/api/jobs/:id")
	// r.GET("/api/jobs/:id/log")

	// create https server with TLS configuration to manages certificates
	s = &Server{
		&http.Server{
			Handler:   r,
			Addr:      port,
			TLSConfig: tlsConfig,
		},
	}
	// srv.Handler = r
	// srv.Addr = ":8080"
	// srv.TLSConfig.Certificates = []tls.Certificate{cert}
	return s
}

// ***TEMP***
//
// Index handles requests to root directory
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Printf("Hello from root")
	msg := map[string]string{"msg": "Hello", "name": "Brady"}

	res, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(res)

}

// get list handler /jobs
func (s *Server) getJobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type response struct{}

	respond()
}

// get 1 job handler /jobs/:id
func (s *Server) getJobStatus(w http.ResponseWriter, r *http.Request, id httprouter.Params) {}

// post 1 job handler /jobs
func (s *Server) addJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

// delete 1 job handler /jobs/:id
func (s *Server) cancelJob(w http.ResponseWriter, r *http.Request, id httprouter.Params) {}

// get 1 job log handler /jobs/:id/log
func (s *Server) getJobLog(w http.ResponseWriter, r *http.Request, id httprouter.Params) {}

func respond() {}
