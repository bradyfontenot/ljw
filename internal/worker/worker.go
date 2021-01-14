package worker

import (
	"errors"
	"strconv"
	"sync"
)

var idCounter int = 0

// Worker is a store and task manager for all Jobs
type Worker struct {
	Jobs   map[int]*Job // key serves as job id
	currID int
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
	wkr.Jobs[wkr.currID] = newJob(cmd)

	// start job
	go func() {
		if err := wkr.Jobs[wkr.currID].start(wkr.currID); err != nil {
			// print error msg srvr side and pass through
			// to handler
			// fmt.Println(err)
			// return -1, err
		}
	}()
	return wkr.currID, nil
}

// ListRunningJobs returns a list of running jobs
func (wkr *Worker) ListRunningJobs() []int {

	var list []int
	for id, job := range wkr.Jobs {
		if running == job.status {
			list = append(list, id)
		}
	}

	return list
}

func (wkr *Worker) StopJob(id string) (bool, error) {
	idInt, _ := strconv.Atoi(id)

	// validate id
	job, ok := wkr.Jobs[idInt]
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
	job, ok := wkr.Jobs[idInt]
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
	job, ok := wkr.Jobs[idInt]
	if !ok {
		return nil, errors.New("invalid id")
	}

	return map[string]string{
		"cmd":    job.cmd,
		"status": job.status,
		"output": job.output,
	}, nil

}
