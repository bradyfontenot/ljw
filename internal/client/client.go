package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	baseURI  = "https://localhost:8080/"
	certFile = "ssl/client.crt"
	keyFile  = "ssl/client.key"
	caFile   = "ssl/ca.crt"
)

// Client implements an http client
type Client struct {
	*http.Client
}

// New creates and returns a new Client
func New() *Client {

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	c := new(Client)

	c = &Client{
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool,
					Certificates: []tls.Certificate{cert},
				},
			},
		},
	}

	return c
}

// ***TEMP****
//
// TestGet is a test function To Be Deleted
func (cl *Client) TestGet() {
	resp, err := cl.Get(baseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("%s", body)
}

// ListJobs requests a list of all running jobs and outputs results to terminal
func (cl *Client) ListJobs() {}

// JobStatus requests the status of job matching id
func (cl *Client) JobStatus(id int) {}

// StartJob posts a request to start a new job
func (cl *Client) StartJob() {}

// StopJob requests to delete job matching id
func (cl *Client) StopJob() {}

// LogJob
func (cl *Client) LogJob() {}
