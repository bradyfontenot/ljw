package worker

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// for debug. temp.
var idCounter int = 0

// Worker is a store and task manager for all jobs
type Worker struct {
	jobs   map[string]*job // key serves as job id
	currID int             // for debug. temp.
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

// ListJobs returns a list of   jobs
func (wkr *Worker) ListJobs() []string {
	wkr.RLock()
	defer wkr.RUnlock()

	var list []string
	for id, job := range wkr.jobs {
		job.RLock()
		defer job.RUnlock()
		list = append(list, id)
	}

	return list
}

// StartJob initializes a new job and makes call to start the proc
// return props of new job
func (wkr *Worker) StartJob(cmd []string) map[string]string {
	wkr.Lock()
	defer wkr.Unlock()

	id := uuid.New().String()

	// for debug. temp. replacing uuid for now
	wkr.currID++
	id = strconv.Itoa(wkr.currID)

	// create new job instance
	wkr.jobs[id] = newJob(cmd)
	job := wkr.jobs[id] // dont need this assignment but easier to read below.

	// start job
	job.start(id)

	job.RLock()
	defer job.RUnlock()
	return map[string]string{
		"id":     id,
		"cmd":    strings.Join(job.cmd, " "),
		"status": job.status,
		"output": job.output,
	}
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

// GetJob returns a map of job props for the matching id
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
