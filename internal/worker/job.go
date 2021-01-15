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
	cmd    []string
	status string
	output string
	pid    int
	mu     sync.Mutex
}

// New creates a new Job
func newJob(cmd []string) *Job {
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

	// init new Command
	cmd := exec.Command(j.cmd[0], j.cmd[1:]...)
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
	// set job status
	j.status = running

	// mutex here
	// store the Pid so stop() can be called later if needed.
	j.pid = cmd.Process.Pid

	// wait for process to end before access buffer below
	cmd.Wait()

	// mutex here
	// store result from process
	stdoutStderr := buf.Bytes()
	j.output = string(stdoutStderr)
	j.status = finished

	fmt.Printf("PID: %d Finished Job # %d \n", j.pid, id)
	fmt.Printf("Output: %s\n", j.output)
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
