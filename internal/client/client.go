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
	baseURI  = "https://localhost:8080"
	certFile = "ssl/client.crt"
	keyFile  = "ssl/client.key"
	caFile   = "ssl/ca.crt"
)

// Client implements an http client
type Client struct {
	*http.Client
}

// New creates and returns a new Client
func New() (*Client, error) {

	// load certs and config TLS for client
	tlsConfig, err := setupTLS()
	if err != nil {
		return nil, err
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
	}, nil
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
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("failed to append certs from pem")
	}

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

	if r.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	var resp response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return err
	}

	fmt.Println("[ALL JOBS]")
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

	fmt.Printf("[JOB ADDED]\n[ID]: \t\t%s\n[COMMAND]: \t%s\n[STATUS]: \t%s\n[OUTPUT]:\n%s\n", resp.ID, resp.Cmd, resp.Status, resp.Output)
	fmt.Print("[OUTPUT END]\n\n")

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

	fmt.Printf("[JOB LOG]\n[ID]: \t\t%s\n[COMMAND]: \t%s\n[STATUS]: \t%s\n[OUTPUT]:\n%s\n", id, resp.Cmd, resp.Status, resp.Output)
	fmt.Print("[OUTPUT END]\n\n")

	return nil
}
