package worker

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

var idCounter int = 0

type Worker struct {
	JobList map[int]*Job // key functions as id
	currID  int
	sync.Mutex
}

// New creates a new Worker
func New() *Worker {
	return &Worker{
		make(map[int]*Job),
		0,
		sync.Mutex{},
	}
}

func (wkr *Worker) StartJob(cmd string) (int, error) {
	defer wkr.Unlock()
	wkr.Lock()
	wkr.currID++

	// create new job instance
	wkr.JobList[wkr.currID] = newJob(cmd)

	// execute job
	if err := wkr.JobList[wkr.currID].start(); err != nil {
		// print error msg srvr side and pass through
		// to handler
		fmt.Println(err)
		return -1, err
	}

	return wkr.currID, nil
}

type RunningJobs struct {
	ID int
	// Cmd string
}

// GetRunningJobs returns a list of running jobs
func (wkr *Worker) ListRunningJobs() ([]RunningJobs, error) {

	// Temp return unfiltered list until jobs implemented.
	var list []RunningJobs
	for k, _ := range wkr.JobList {
		job := RunningJobs{k}
		list = append(list, job)
	}

	return list, nil
}

func (wkr *Worker) StopJob(id string) (bool, error) {
	idInt, _ := strconv.Atoi(id)

	// validate id
	job, ok := wkr.JobList[idInt]
	if !ok {
		return false, errors.New("id does not exist")
	}
	// stop job
	if err := job.stop(); err != nil {
		return false, nil
	}

	// going away. will be set in Job.stop()

	return true, nil
}

// func (wkr *Worker) GetJobLog(id string) (Job, error) {
// 	idInt, _ := strconv.Atoi(id)

// 	return *wkr.JobList[idInt], nil
// }

func (wkr *Worker) GetJobStatus(id string) (map[string]string, error) {
	// TODO: handle nonexistent/invalid ID
	idInt, _ := strconv.Atoi(id)

	// validate id
	if _, ok := wkr.JobList[idInt]; !ok {
		return nil, errors.New("id does not exist")
	}

	m := map[string]string{
		"status": wkr.JobList[idInt].status,
		"output": wkr.JobList[idInt].output,
	}

	return m, nil
}

func (wkr *Worker) GetJobLog(id string) (map[string]string, error) {
	idInt, _ := strconv.Atoi(id)

	// validate id exists
	job, ok := wkr.JobList[idInt]
	if !ok {
		return nil, errors.New("id does not exist")
	}

	log := map[string]string{
		"cmd":    job.cmd,
		"status": job.status,
		"output": job.output,
	}

	return log, nil
}
