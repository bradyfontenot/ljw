package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bradyfontenot/ljw/internal/worker"
)

const (
	addr     = "localhost:8080"
	certFile = "ssl/server.crt"
	keyFile  = "ssl/server.key"
	caFile   = "ssl/ca.crt"
)

// Server implements http server and uses a worker to execute tasks.
type Server struct {
	*http.Server
	worker *worker.Worker
}

// New creates and returns a new server.
func New(wkr *worker.Worker) (*Server, error) {

	// load certs and config TLS for server
	tlsConfig, err := setupTLS()
	if err != nil {
		return nil, err
	}

	var s Server
	s = Server{
		&http.Server{
			Addr:    addr,
			Handler: s.router(),
			// Generic timeout. Could use header timeout if you want
			// to set specific read timeout for each handler
			ReadTimeout:  time.Duration(30 * time.Second),
			WriteTimeout: time.Duration(30 * time.Second),
			TLSConfig:    tlsConfig,
		},
		wkr,
	}

	return &s, nil
}

// setupTLS sets up Authentication and builds tlsConfig for the server.
func setupTLS() (*tls.Config, error) {

	// load certificate authority file
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	// create pool for accepted certificate authorities and add ca.
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("failed to append certs from pem")
	}

	// load certificate and private key files
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		Certificates: []tls.Certificate{cert},
	}, nil
}
