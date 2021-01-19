package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	// "github.com/bradyfontenot/ljw/internal/server"
	"github.com/bradyfontenot/ljw/internal/worker"
	"github.com/stretchr/testify/assert"
)

func TestStartJob(t *testing.T) {
	// create server and populate job worker
	srv := New(worker.New())

	// create command
	cmd := map[string][]string{"Cmd": []string{"echo", "Hello", "World"}}
	// handle err?
	reqBody, _ := json.Marshal(cmd)

	t.Run("starting a job", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/jobs"), bytes.NewBuffer(reqBody))
		resp := httptest.NewRecorder()

		srv.Handler.ServeHTTP(resp, req)
		respResult := resp.Result()

		// test response status code
		assert.Equal(t, http.StatusCreated, respResult.StatusCode, "status code does not match")

		// test json format
		actualJSON, _ := ioutil.ReadAll(respResult.Body)
		expectedJSON := `{"id":"1","cmd":"echo Hello World", "status":"FINISHED", "output":"Hello World\n"}`
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
			Status: "FINISHED",
			Output: "Hello World\n",
		}

		assert.Equal(t, expected, actual)
	})
}

// decode json response into map
func decodeRes(res *httptest.ResponseRecorder) map[string]string {
	var got map[string]string
	json.NewDecoder(res.Body).Decode(&got)

	return got
}

func TestStopJob(t *testing.T) {
	// create server and populate job worker
	srv := New(worker.New())

	// create command and start job to be canceled
	// id will be 1
	cmd := []string{"sleep", "30"}
	srv.worker.StartJob(cmd)
	time.Sleep(100 * time.Millisecond)

	t.Run("successful stop request", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/jobs/1"), nil)
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
}

func TestGetJob(t *testing.T) {
	// create server and populate job worker
	srv := New(worker.New())
	// create command and start job.
	id := "1"
	cmd := []string{"echo", "Hello Teleport"}
	srv.worker.StartJob(cmd)

	t.Run("successful job request", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/jobs/%s", id), nil)
		resp := httptest.NewRecorder()

		srv.Handler.ServeHTTP(resp, req)

		// test response status code
		respResult := resp.Result()
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
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/jobs/%s", id), nil)
		resp := httptest.NewRecorder()

		srv.Handler.ServeHTTP(resp, req)

		// test response status code
		respResult := resp.Result()
		assert.Equal(t, http.StatusNotFound, respResult.StatusCode, "status code does not match")

		// check error message
		actual, _ := ioutil.ReadAll(respResult.Body)
		expected := "invalid id\n"

		assert.Equal(t, expected, string(actual), "error message should be 'invalid id'")

	})
}
