// Package server is an https server
package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bradyfontenot/ljw/internal/worker"
)

const (
	port     = ":8080"
	certFile = "ssl/server.crt"
	keyFile  = "ssl/server.key"
	caFile   = "ssl/ca.crt"
)

// Server
type Server struct {
	*http.Server
	worker *worker.Worker
}

// New creates a new server.
func New(wkr *worker.Worker) *Server {

	// load certs and config TLS for server
	tlsConfig, err := setupTLS()
	if err != nil {
		log.Fatal(err)
	}

	var s Server
	s = Server{
		&http.Server{
			Addr:    port,
			Handler: s.router(),
			// Generic timeout. Could use header timeout if you want to set specific read timeout for each handler
			ReadTimeout:  time.Duration(30 * time.Second),
			WriteTimeout: time.Duration(30 * time.Second),
			TLSConfig:    tlsConfig,
		},
		wkr,
	}

	return &s
}

// buildTLSConfig setups Authentication and builds tlsConfig for the server.
func setupTLS() (*tls.Config, error) {

	// load certificate authority file
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return &tls.Config{}, err
	}

	// create pool for accepted certificate authorities and add ca.
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// load certificate and private key files
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return &tls.Config{}, err
	}

	return &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		Certificates: []tls.Certificate{cert},
	}, nil
}
