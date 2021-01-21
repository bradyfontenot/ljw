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
*/

package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

	return &Client{
		&http.Client{
			Timeout: time.Duration(30 * time.Second),
		},
	}
}

// SetupTLS sets up Authentication and builds tlsConfig for the client
func (cl *Client) SetupTLS() error {

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return errors.New("failed to append certs from pem")
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	tls := tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}

	tr := &http.Transport{
		TLSClientConfig:     &tls,
		TLSHandshakeTimeout: time.Duration(15 * time.Second),
	}

	cl.Transport = tr
	return nil
}

// ListJobs requests a list of all jobs and outputs id and status
func (cl *Client) ListJobs() error {
	type response struct {
		IDList []string `json:"idList"`
	}

	r, err := cl.Get(baseURI + "/api/jobs")
	if err != nil {
		return err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var resp response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return err
	}

	// TODO:
	//	Printing output here for simplicity.
	//	In mvp/prod we should return (response, err) for flexibility
	//	to format/handle data on frontend.
	fmt.Println("[ ALL JOBS]")
	for _, v := range resp.IDList {
		fmt.Println(" -ID:", v)
	}

	return nil
}

// StartJob posts a request to start a new job
func (cl *Client) StartJob(cmd []string) error {
	type request struct {
		Cmd []string `json:"cmd"`
	}

	type response struct {
		ID     string `json:"id"`
		Cmd    string `json:"cmd"`
		Status string `json:"status"`
		Output string `json:"output"`
	}

	msg := request{cmd}
	reqBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	r, err := cl.Post(baseURI+"/api/jobs", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusCreated {
		return errors.New(string(body))
	}

	var resp response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return err
	}

	fmt.Printf("[JOB ADDED]\n -ID: %s\n -Status: %s\n", resp.ID, resp.Status)
	if resp.Status == "FINISHED" || resp.Status == "FAILED" {
		fmt.Printf(" -Output:\n%s", resp.Output)
		fmt.Print("[END]\n\n")
	}
	return nil
}

// JobStatus requests the status of job matching id
func (cl *Client) JobStatus(id string) error {
	type response struct {
		Status string `json:"status"`
	}

	r, err := cl.Get(baseURI + "/api/jobs/" + id)
	if err != nil {
		return err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	var resp response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return err
	}

	fmt.Printf("[JOB STATUS] => %s \n", resp.Status)

	return nil
}

// StopJob requests to delete job matching id
func (cl *Client) StopJob(id string) error {
	type response struct {
		Success bool `json:"success"`
	}

	req, err := http.NewRequest("DELETE", baseURI+"/api/jobs/"+id, nil)
	if err != nil {
		return err
	}

	r, err := cl.Do(req)
	if err != nil {
		return err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	var resp response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return err
	}

	fmt.Printf("[JOB STOPPED] => %s \n", strings.ToUpper(strconv.FormatBool(resp.Success)))

	return nil
}

// GetJobLog ....
func (cl *Client) GetJobLog(id string) error {
	type response struct {
		Cmd    string `json:"cmd"`
		Status string `json:"status"`
		Output string `json:"output"`
	}

	r, err := cl.Get(baseURI + "/api/jobs/" + id + "/log")
	if err != nil {
		return err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	var resp response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return err
	}

	fmt.Printf("[JOB LOG]\n -ID: %s\n -Command: %s\n -Status: %s\n -Output:\n%s\n", id, resp.Cmd, resp.Status, resp.Output)
	fmt.Print("[END]\n\n")

	return nil
}
