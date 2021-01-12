package server

import (
	"encoding/json"
	"net/http"

	"github.com/bradyfontenot/ljw/internal/worker"
	"github.com/julienschmidt/httprouter"
)

// Response is ...
type Response struct {
	Msg     string               `json:"msg,omitempty"`
	Error   string               `json:"error,omitempty"`
	ID      int                  `json:"id,omitempty"`
	Status  string               `json:"status,omitempty"`
	Cmd     string               `json:"cmd,omitempty"`
	Output  string               `json:"output,omitempty"`
	JobList []worker.RunningJobs `json:"jobList,omitempty"`
}

func (s *Server) router() *httprouter.Router {

	r := httprouter.New()

	r.GET("/api/jobs", s.listJobs)
	r.GET("/api/jobs/:id", s.getJobStatus)
	r.POST("/api/jobs", s.startJob)
	r.DELETE("/api/jobs/:id", s.stopJob)
	r.GET("/api/jobs/:id/log", s.getJobLog)

	return r
}

// createJob calls the worker to create and start a new job to execute the linux command
// in the client's request msg. If job is successfully added, the new job's id will be sent
// to the client as a response along w/ a 201 statuscode.
func (s *Server) startJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type request struct {
		Cmd string
	}

	// decode request msg
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// pass cmd to worker to build new job and receive id of new job
	jobID, err := s.worker.StartJob(req.Cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// set header properties
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// build response msg & send
	resp := Response{ID: jobID}
	sendResp(w, resp)
}

// getRunningJobs gets list of jobs currently in process
//
// returns array of objects with job id and cmd properties
func (s *Server) listJobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// get list of running jobs
	jobList, err := s.worker.ListJobs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// set header properties
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// build response msg & send
	resp := Response{
		JobList: jobList, // [{id, cmd}]
	}
	sendResp(w, resp)
}

// getJobStatus returns status of job matching id. If status is "complete", the output will
// also be included in response
func (s *Server) getJobStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	status, err := s.worker.GetJobStatus(p.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// set header properties
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// build response msg & send
	resp := Response{
		Status: status["status"],
		Output: status["output"],
	}
	sendResp(w, resp)
}

// stopJob stops job matching id if currently running. Does nothing if job not running
func (s *Server) stopJob(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	msg, err := s.worker.StopJob(p.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// set header properties
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// build response msg & send
	resp := Response{
		Msg: msg,
	}
	sendResp(w, resp)
}

// getJobLog returns log for job matching id
func (s *Server) getJobLog(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log, err := s.worker.GetJobLog(p.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// set header properties
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// build response msg & send
	resp := Response{
		Cmd:    log["cmd"],
		Status: log["status"],
		Output: log["output"],
	}
	sendResp(w, resp)
}

// helper function for marshalling json & sending response
func sendResp(w http.ResponseWriter, msg Response) {
	resp, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// fmt.Printf("%s", resp)
	w.Write(resp)
}
