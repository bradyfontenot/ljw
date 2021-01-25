/*
/	router.go defines all the routes and handlers
/	I broke this out to keep server.go from getting cluttered.
/
/	For a larger project you'd probably want to break this out
/	into its own package w/ separate files(or even separate packages)
/	for each major endpoint and it's handlers depending on scale of project
*/

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Response is a catch all response struct
//
// I decided to use this instead of defining a struct for each response separately so that
// I could pass a single type into the sendResp() helper that marshals the json and writes the request.
// I wanted to avoid some repetitive boilerplate in the handlers.
type Response struct {
	Success bool     `json:"success,omitempty"` // successful operation
	ID      string   `json:"id,omitempty"`      // job ID
	Status  string   `json:"status,omitempty"`  // job status
	Cmd     string   `json:"cmd,omitempty"`     // job command
	Output  string   `json:"output,omitempty"`  // job output
	IDList  []string `json:"idList,omitempty"`  // list of job ID's
}

// router creates handler and defines the routes.
func (s *Server) router() *httprouter.Router {

	r := httprouter.New()

	r.GET("/api/jobs", s.listJobs)
	r.POST("/api/jobs", s.startJob)
	r.GET("/api/jobs/:id", s.getJob)
	r.DELETE("/api/jobs/:id", s.stopJob)
	r.GET("/api/jobs/:id/log", s.getJob)

	return r
}

// listJobs retrieves list of ids for jobs currently in process
func (s *Server) listJobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// get list of jobs
	idList := s.worker.ListJobs()

	// set header properties
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// build response msg & send
	resp := Response{
		IDList: idList,
	}
	sendResp(w, resp)
}

// startJob starts a new job and returns new job id if successful
func (s *Server) startJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type request struct {
		Cmd []string
	}
	// decode request msg
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// pass cmd to worker to build new job and receive job props
	job := s.worker.StartJob(req.Cmd)

	// set header properties
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// build response msg & send
	resp := Response{
		ID:     job["id"],
		Cmd:    job["cmd"],
		Status: job["status"],
		Output: job["output"],
	}
	sendResp(w, resp)
}

// stopJob stops job if it is currently running.
// returns a boolean to confirm if job was canceled or not
func (s *Server) stopJob(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	result, err := s.worker.StopJob(p.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// set header properties
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// build response msg & send
	//
	// probably should include reason if job isn't cancelled b/c
	// it's already complete so client knows why job wasn't cancelled.
	resp := Response{
		Success: result,
	}
	sendResp(w, resp)
}

// getJob returns job matching id
// called by client funcs: JobLog() & JobStatus()
func (s *Server) getJob(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	job, err := s.worker.GetJob(p.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// set header properties
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// build response msg & send
	resp := Response{
		ID:     job["id"],
		Cmd:    job["cmd"],
		Status: job["status"],
		Output: job["output"],
	}
	sendResp(w, resp)
}

// helper function for marshalling json & sending response
func sendResp(w http.ResponseWriter, msg Response) {
	resp, err := json.Marshal(msg)
	if err != nil {
		e := fmt.Errorf("could not marshall json. error: %w", err)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}
