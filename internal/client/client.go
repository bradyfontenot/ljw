package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		fmt.Println(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Println(err)
	}
	c := new(Client)

	c = &Client{
		&http.Client{
			Timeout: time.Duration(30 * time.Second),
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

// GetJobs requests a list of all running jobs and outputs results to terminal
func (cl *Client) GetJobs() {
	type response struct {
		JobList []struct {
			ID  int    `json:"id"`
			Cmd string `json:"cmd"`
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

	r, err := cl.Get(baseURI + "/api/jobs/" + id)
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
