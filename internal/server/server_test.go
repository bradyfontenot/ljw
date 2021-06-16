// TODO:
//	some of this code is repetitive and should be put into
//	helper functions.

package server

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/bradyfontenot/ljw/internal/client"
	"github.com/bradyfontenot/ljw/internal/worker"
	"github.com/stretchr/testify/assert"
)

func init() {
	// set working directory to project root
	// needed so that cert files can be found by tests
	// since file paths are hardcoded
	if err := os.Chdir("../.."); err != nil {
		log.Fatal(err)
	}
}

func TestStartJob(t *testing.T) {
	// create server and populate job worker
	srv, err := New(worker.New())
	if err != nil {
		log.Fatal(err)
	}
	cmd := map[string][]string{"Cmd": []string{"echo", "Hello", "World"}}
	reqBody, err := json.Marshal(cmd)
	if err != nil {
		log.Fatal(err)
	}

	t.Run("starting a job", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/jobs"), bytes.NewBuffer(reqBody))
		resp := httptest.NewRecorder()

		srv.Handler.ServeHTTP(resp, req)
		respResult := resp.Result()

		// test response status code
		assert.Equal(t, http.StatusCreated, respResult.StatusCode, "status code does not match")

		// test json format
		actualJSON, _ := ioutil.ReadAll(respResult.Body)
		expectedJSON := `{"id":"1","cmd":"echo Hello World", "status":"RUNNING"}`
		assert.JSONEq(t, expectedJSON, string(actualJSON), "json does not match")

		// test data structure
		type response struct {
			ID     string
			Cmd    string
			Status string
			Output string
		}

		var actual response
		json.Unmarshal(actualJSON, &actual)
		expected := response{
			ID:     "1",
			Cmd:    "echo Hello World",
			Status: "RUNNING",
			Output: "",
		}

		assert.Equal(t, expected, actual)
	})
}

func TestStopJob(t *testing.T) {

	// create server and populate worker with a job
	srv, err := New(worker.New())
	if err != nil {
		log.Fatal(err)
	}
	id := "1"
	cmd := []string{"sleep", "2"}
	srv.worker.StartJob(cmd)

	t.Run("successful stop request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/jobs/%s", id), nil)
		resp := httptest.NewRecorder()

		srv.Handler.ServeHTTP(resp, req)

		// test response status code
		respResult := resp.Result()
		assert.Equal(t, http.StatusOK, respResult.StatusCode, "status code does not match")

		// test json format
		actualJSON, _ := ioutil.ReadAll(respResult.Body)
		expectedJSON := `{"success":true}`
		assert.JSONEq(t, expectedJSON, string(actualJSON), "json does not match")

		// test data structure
		type response struct {
			Success bool
		}
		var actual response
		json.Unmarshal(actualJSON, &actual)
		expected := response{
			Success: true,
		}

		assert.Equal(t, expected, actual, "response structure does not match")
	})

	t.Run("stop request with nonexistent id", func(t *testing.T) {
		id = "5"
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/jobs/%s", id), nil)
		resp := httptest.NewRecorder()

		srv.Handler.ServeHTTP(resp, req)

		// test response status code
		respResult := resp.Result()
		assert.Equal(t, http.StatusNotFound, respResult.StatusCode, "status code does not match")

		// check error message
		actual, _ := ioutil.ReadAll(respResult.Body)
		expected := fmt.Sprintf("%s is not a valid id\n", id)

		assert.Equal(t, expected, string(actual))
	})
}

func TestGetJob(t *testing.T) {
	// create server and populate worker w/ a job
	srv, err := New(worker.New())
	if err != nil {
		log.Fatal(err)
	}

	id := "1"
	cmd := []string{"echo", "Hello Teleport"}
	srv.worker.StartJob(cmd)
	// give command a little time to finish before checking log for output
	time.Sleep(25 * time.Millisecond)

	t.Run("successful job request using valid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/jobs/%s", id), nil)
		resp := httptest.NewRecorder()

		srv.Handler.ServeHTTP(resp, req)
		respResult := resp.Result()

		// test response status code
		assert.Equal(t, http.StatusOK, respResult.StatusCode, "status code does not match")

		// test json format
		actualJSON, _ := ioutil.ReadAll(respResult.Body)
		expectedJSON := `{"id":"1", "cmd":"echo Hello Teleport", "status":"FINISHED", "output":"Hello Teleport\n"}`

		assert.JSONEq(t, expectedJSON, string(actualJSON), "json does not match")

		// test data structure
		type response struct {
			ID     string
			Cmd    string
			Status string
			Output string
		}
		var actual response
		json.Unmarshal(actualJSON, &actual)
		expected := response{
			ID:     "1",
			Cmd:    "echo Hello Teleport",
			Status: "FINISHED",
			Output: "Hello Teleport\n",
		}

		assert.Equal(t, expected, actual, "response structure does not match")
	})

	t.Run("invalid job request using noexistent id", func(t *testing.T) {
		id := "2"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/jobs/%s", id), nil)
		resp := httptest.NewRecorder()

		srv.Handler.ServeHTTP(resp, req)

		// test response status code
		respResult := resp.Result()
		assert.Equal(t, http.StatusNotFound, respResult.StatusCode, "status code does not match")

		// check error message
		actual, _ := ioutil.ReadAll(respResult.Body)
		expected := fmt.Sprintf("%s is not a valid id\n", id)

		assert.Equal(t, expected, string(actual))

	})
}

func TestClientAuthentication(t *testing.T) {

	t.Run("test valid client connection is accepted", func(t *testing.T) {

		cl, err := client.New()
		if err != nil {
			log.Fatal(err)
		}

		// creat a new server and assign it's
		// configuration to the httptest's TLS
		// server for testing.
		srv, err := New(worker.New())
		if err != nil {
			log.Fatal(err)
		}

		s := httptest.NewUnstartedServer(srv.Handler)
		s.TLS = srv.TLSConfig
		s.StartTLS()

		res, err := cl.Get(s.URL + "/api/jobs")
		if err != nil {
			log.Fatal(err)
		}

		actual := res.StatusCode
		assert.Equal(t, http.StatusOK, actual)

		s.Close()
	})

	t.Run("test invalid client connection is rejected", func(t *testing.T) {
		// init client and reconfigure with invalid certificates.
		cl, err := client.New()
		if err != nil {
			log.Fatal(err)
		}
		tls, err := invalidClientTLS()
		if err != nil {
			log.Fatal(err)
		}
		tr := &http.Transport{
			TLSClientConfig:     tls,
			TLSHandshakeTimeout: time.Duration(15 * time.Second),
		}
		cl.Transport = tr

		// configure  and run server
		// assign handler and tlsconfig from app's server to TLS test server
		srv, err := New(worker.New())
		if err != nil {
			log.Fatal(err)
		}
		s := httptest.NewUnstartedServer(srv.Handler)
		s.TLS = srv.TLSConfig
		s.StartTLS()

		// build client request and check for failure
		_, err = cl.Get(s.URL + "/api/jobs")
		if err != nil {
			// pass through to prevent compile error exit code in test result
		}
		assert.Contains(t, err.Error(), "x509: certificate signed by unknown authority")

	})

}

// builds a tls config with invalid certifcates that is assigned
// to previously created client
func invalidClientTLS() (*tls.Config, error) {
	certFile := "ssl/invalid_test_ssl/invalid_client.crt"
	keyFile := "ssl/invalid_test_ssl/invalid_client.key"
	caFile := "ssl/invalid_test_ssl/invalid_ca.crt"

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
