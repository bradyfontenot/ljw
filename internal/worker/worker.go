package worker

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

var idCounter int = 0

// Worker is a store and task manager for all jobs
type Worker struct {
	jobs   map[int]*job // key serves as job id
	currID int
	sync.Mutex
}

// New creates a new Worker
func New() *Worker {
	return &Worker{
		make(map[int]*job),
		0,
		sync.Mutex{},
	}
}

func (wkr *Worker) StartJob(cmd []string) (int, error) {
	wkr.currID++
	id := wkr.currID
	// create new job instance
	wkr.jobs[id] = newJob(cmd)

	// start job
	go func() {
		if err := wkr.jobs[wkr.currID].start(wkr.currID); err != nil {
			// print error msg srvr side and pass through
			// to handler
			// fmt.Println(err)
			// return -1, err
		}
	}()
	return id, nil
}

// ListRunningJobs returns a list of running jobs
func (wkr *Worker) ListRunningJobs() []int {
	var list []int
	for id, job := range wkr.jobs {
		if running == job.status {
			list = append(list, id)
		}
	}

	return list
}

func (wkr *Worker) StopJob(id string) (bool, error) {
	idInt, _ := strconv.Atoi(id)

	// validate id
	job, ok := wkr.jobs[idInt]
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

func (wkr *Worker) GetJobStatus(id string) (map[string]string, error) {
	idInt, _ := strconv.Atoi(id)

	// validate id
	job, ok := wkr.jobs[idInt]
	if !ok {
		return nil, errors.New("invalid id")
	}

	return map[string]string{
		"status": job.status,
		"output": job.output,
	}, nil
}

func (wkr *Worker) GetJob(id string) (map[string]string, error) {
	idInt, _ := strconv.Atoi(id)

	// validate id
	job, ok := wkr.jobs[idInt]
	if !ok {
		return nil, errors.New("invalid id")
	}

	return map[string]string{
		"cmd":    strings.Join(job.cmd, " "),
		"status": job.status,
		"output": job.output,
	}, nil

}
