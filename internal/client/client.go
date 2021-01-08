package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	baseURI = "http://localhost:8080/"
)

// Client implements an http client
type Client struct {
	http.Client
	// tr http.Transport
}

// New creates and returns a new Client
func New() *Client {
	client := new(Client)

	return client
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
