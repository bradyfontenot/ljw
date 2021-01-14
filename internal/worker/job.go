package worker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

const (
	queued   = "QUEUED"
	running  = "RUNNING"
	finished = "FINISHED"
	canceled = "CANCELED"
	failed   = "FAILED"
)

// Job is one linux job
type Job struct {
	cmd    string
	status string
	output string
	pid    int
	mu     sync.Mutex
}

// New creates a new Job
func newJob(cmd string) *Job {
	return &Job{
		cmd:    cmd,
		status: running,
		mu:     sync.Mutex{},
	}
}

// TODO:
//	* handle invalid cmd(don't kill server. just set status to Failed)
//	* Run in go routine ya?

func (j *Job) start(id int) error {
	// set job status
	j.status = running

	// init new Command
	cmd := exec.Command("sleep", "5")

	// set stdout and stderr to write to same buffer
	// mimics CombinedOutput() from os/exec library
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	// run command
	if err := cmd.Start(); err != nil {
		// mutex here probably.
		j.status = failed
		j.output = err.Error()

		return fmt.Errorf("error starting command: %w", err)
	}

	// mutex here
	// store the Pid so stop() can be called later if needed.
	j.pid = cmd.Process.Pid

	// get stdout/stderr from buffer
	var stdoutStderr []byte
	buf.Write(stdoutStderr)

	cmd.Wait()

	// mutex here
	// store result from process
	j.output = string(stdoutStderr)
	j.status = finished

	fmt.Printf("%d Finished Job # %d: %s\n", j.pid, id, stdoutStderr)

	return nil
}

func (j *Job) stop() (bool, error) {
	if j.status != running {
		return false, nil
	}

	proc, err := os.FindProcess(j.pid)
	if err != nil {
		return false, fmt.Errorf("could not find process. error: %v", err)
	}
	if err = proc.Kill(); err != nil {
		return false, fmt.Errorf("could not kill process. error: %v", err)
	}

	j.status = canceled

	return true, nil
}
