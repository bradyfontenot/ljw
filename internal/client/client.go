/*
/	Currently contains client constructor, auth function and
/	request functions.  Similar to server, auth could be
/	broken out into separate package if warranted by complexity.
/
/	Again, could also break out the requests if desired for a larger project
/	with more than one endpoint and have a package for each.
/	I've not done an api client in Go before but, in the past w/ elixir
/	I've had a module for the client w/ generic methods for each request type
/	and then separate modules for each endpoint w/ relevant methods.
/
/	So you could do something like set a baseURL w/ some middleware and then
/ 	just pass an endpoint to the client's generic request method instead of
/	passing the baseURL and hardcoding the endpoint everytime like below.
/
/	cert/key/ca files stored in repo w/ paths hardcoded
/	but should be accessed using environment variables
/	or other method to keep hidden/secure.
/
/	TODO:
/		* finish cli
/		* write some helper funcs to format/print output.
*/

package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

	// load certs and config TLS for client
	tlsConfig, err := setupTLS()
	if err != nil {
		log.Fatal(err)
	}

	tr := &http.Transport{
		TLSClientConfig:     tlsConfig,
		TLSHandshakeTimeout: time.Duration(15 * time.Second),
	}

	return &Client{
		&http.Client{
			Timeout:   time.Duration(30 * time.Second),
			Transport: tr,
		},
	}
}

// setupTLS sets up Authentication and builds tlsConfig for the client
func setupTLS() (*tls.Config, error) {

	// load certificate authority file
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	// create pool for accepted certificate authorities and add ca.
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// load certificate and private key files
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}, nil
}

// ListJobs requests a list of all jobs and outputs id and status
//
// **** QUESTION:  Should this only be running jobs or all jobs? *******
func (cl *Client) ListRunningJobs() {
	type response struct {
		JobList []struct {
			ID int `json:"id"`
		} `json:"jobList"`
	}

	r, err := cl.Get(baseURI + "/api/jobs")
	if err != nil {
		fmt.Println(err)
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	// extract response msg
	var resp response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		fmt.Println(err)
	}

	// TODO: Format Output
	fmt.Printf("Running Jobs: %+v \n", resp)
}

// StartJob posts a request to start a new job
func (cl *Client) StartJob() {
	type request struct {
		Cmd string `json:"cmd"`
	}

	type response struct {
		ID int `json:"id"`
	}

	// temporary cmd until cli implemented.
	msg := request{"test command"}
	reqBody, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}

	// send request & capture response
	r, err := cl.Post(baseURI+"/api/jobs", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
	}

	defer r.Body.Close()
	// TODO: handle error
	body, err := ioutil.ReadAll(r.Body)

	// extract response msg
	var resp response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		fmt.Println(err)
	}

	// TODO: Format Output
	fmt.Printf("Job Added: %+v \n", resp)
}

// JobStatus requests the status of job matching id
func (cl *Client) JobStatus(id string) {
	// send request & capture response
	r, err := cl.Get(baseURI + "/api/jobs/" + id)
	if err != nil {
		fmt.Println(err)
	}

	type response struct {
		Status string `json:"status"`
		Output string `json:"output,omitempty"`
	}

	defer r.Body.Close()
	// TODO: handle error
	body, err := ioutil.ReadAll(r.Body)

	// extract response msg
	var resp response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		fmt.Println(err)
	}

	// TODO: Format Output
	fmt.Printf("Job Status: %+v \n", resp)
}

// StopJob requests to delete job matching id
func (cl *Client) StopJob(id string) {
	// build DELETE request
	req, err := http.NewRequest("DELETE", baseURI+"/api/jobs/"+id, nil)
	if err != nil {
		fmt.Println(err)
	}

	// capture response
	r, err := cl.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	type response struct {
		Success bool `json:"success"`
	}

	defer r.Body.Close()
	// TODO: handle error
	body, err := ioutil.ReadAll(r.Body)

	// extract response msg
	var resp response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		fmt.Println(err)
	}

	// TODO: Format Output
	fmt.Printf("Job Canceled?: %+v \n", resp)
}

// GetJobLog ....
func (cl *Client) GetJobLog(id string) {

	// send request & capture response
	r, err := cl.Get(baseURI + "/api/jobs/" + id + "/log")
	if err != nil {
		fmt.Println(err)
	}

	type response struct {
		Cmd    string `json:"cmd"`
		Status string `json:"status"`
		Output string `json:"output"`
	}

	defer r.Body.Close()
	// TODO: handle error
	body, err := ioutil.ReadAll(r.Body)

	// extract response msg
	var resp response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		fmt.Println(err)
	}

	// TODO: Format Output
	fmt.Printf("Job Log: %+v \n", resp)
}
