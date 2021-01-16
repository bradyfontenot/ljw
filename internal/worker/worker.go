package worker

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

var idCounter int = 0

// Worker is a store and task manager for all jobs
type Worker struct {
	jobs   map[string]*job // key serves as job id
	currID int
	*sync.RWMutex
}

// New creates a new Worker
func New() *Worker {
	return &Worker{
		make(map[string]*job),
		0,
		&sync.RWMutex{},
	}
}

// ListRunningJobs returns a list of running jobs
func (wkr *Worker) ListRunningJobs() []string {
	wkr.RLock()
	defer wkr.RUnlock()

	var list []string
	for id, job := range wkr.jobs {
		job.RLock()
		defer job.RUnlock()
		// if running == job.status {
		list = append(list, id)
		// }
	}

	return list
}

// StartJob initializes a new job and makes call to start the proc
func (wkr *Worker) StartJob(cmd []string) (string, error) {
	// errChan := make(chan error)
	wkr.Lock()
	defer wkr.Unlock()
	wkr.currID++
	id := strconv.Itoa(wkr.currID)
	id = uuid.New().String()
	// create new job instance
	wkr.jobs[id] = newJob(cmd)

	// start job
	err := wkr.jobs[id].start(id)
	if err != nil {
		return "", err
	}

	// return map[string]string{
	// 	"id":     id,
	// 	"cmd":    "x", //strings.Join(wkr.jobs[id].cmd, " "),
	// 	"status": wkr.jobs[id].status,
	// 	"output": wkr.jobs[id].output,
	// }
	return id, nil
}

// StopJob will cancel job if still running or queued
func (wkr *Worker) StopJob(id string) (bool, error) {
	wkr.Lock()
	defer wkr.Unlock()

	// validate id
	job, ok := wkr.jobs[id]
	if !ok {
		return false, errors.New("invalid id")
	}
	// stop job
	result, err := job.stop()
	if err != nil {
		return false, err
	}

	return result, nil
}

// GetJobStatus retrieves status of job matching id arg
func (wkr *Worker) GetJobStatus(id string) (map[string]string, error) {
	wkr.RLock()
	defer wkr.RUnlock()

	// validate id
	job, ok := wkr.jobs[id]
	if !ok {
		return nil, errors.New("invalid id")
	}
	job.RLock()
	defer job.RUnlock()
	return map[string]string{
		"status": job.status,
		"output": job.output,
	}, nil
}

// GetJob returns a job struct matching id arg
func (wkr *Worker) GetJob(id string) (map[string]string, error) {
	wkr.RLock()
	defer wkr.RUnlock()

	// validate id
	job, ok := wkr.jobs[id]
	if !ok {
		return nil, errors.New("invalid id")
	}
	job.RLock()
	defer job.RUnlock()
	return map[string]string{
		"id":     id,
		"cmd":    strings.Join(job.cmd, " "),
		"status": job.status,
		"output": job.output,
	}, nil
}
