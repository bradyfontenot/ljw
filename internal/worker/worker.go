package worker

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// temp
var idCounter int = 0

// Worker is a store and task manager for all jobs
type Worker struct {
	// key serves as job id
	jobs map[string]*job
	// should be replaced by UUID in production
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

// ListJobs returns a list of jobs
func (wkr *Worker) ListJobs() []string {
	wkr.RLock()
	defer wkr.RUnlock()

	var list []string
	for id := range wkr.jobs {
		list = append(list, id)
	}
	sort.Slice(list, func(i, j int) bool {
		li, _ := strconv.Atoi(list[i])
		lj, _ := strconv.Atoi(list[j])

		return li < lj
	})
	return list
}

// StartJob initializes a new job and makes call to start the proc
// return props of new job
func (wkr *Worker) StartJob(cmd []string) map[string]string {
	wkr.Lock()
	defer wkr.Unlock()

	id := uuid.New().String()

	// temp. replace w/ UUID in prod
	wkr.currID++
	id = strconv.Itoa(wkr.currID)

	wkr.jobs[id] = newJob(cmd)
	job := wkr.jobs[id]

	job.start()
	return map[string]string{
		"id":     id,
		"cmd":    strings.Join(job.Cmd(), " "),
		"status": job.Status(),
		"output": strings.Join(job.Output(), " "),
	}
}

// StopJob will cancel job if still running or queued
func (wkr *Worker) StopJob(id string) (bool, error) {
	wkr.Lock()
	defer wkr.Unlock()

	job, ok := wkr.jobs[id]
	if !ok {
		return false, errors.New(id + " is not a valid id")
	}

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

	job, ok := wkr.jobs[id]
	if !ok {
		return nil, errors.New(id + " is not a valid id")
	}

	return map[string]string{
		"id":     id,
		"cmd":    strings.Join(job.Cmd(), " "),
		"status": job.Status(),
		"output": strings.Join(job.Output(), " "),
	}, nil
}
