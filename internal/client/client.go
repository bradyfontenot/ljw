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

	tlsConfig, err := setupTLS()
	if err != nil {
		log.Fatal(err)
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	c := &Client{
		&http.Client{
			Timeout:   time.Duration(30 * time.Second),
			Transport: tr,
		},
	}

	return c
}

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
func (cl *Client) ListJobs() {
	type response struct {
		JobList []struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
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

	msg := request{"test command"}
	reqBody, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}

	// capture response
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
	req, err := http.NewRequest("DELETE", baseURI+"/api/jobs/"+id, nil)
	if err != nil {
		fmt.Println(err)
	}

	r, err := cl.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	type response struct {
		Msg string `json:"msg"`
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
	type response struct {
		Cmd    string `json:"cmd"`
		Status string `json:"status"`
		Output string `json:"output"`
	}

	r, err := cl.Get(baseURI + "/api/jobs/" + id + "/log")
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
	fmt.Printf("Job Log: %+v \n", resp)
}
