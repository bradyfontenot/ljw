package worker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

// status values
const (
	queued   = "QUEUED"
	running  = "RUNNING"
	finished = "FINISHED"
	canceled = "CANCELED"
	failed   = "FAILED"
)

// job is one linux job
type job struct {
	cmd    []string
	status string
	output string
	pid    int
	sync.RWMutex
}

// New creates a new job
func newJob(cmd []string) *job {
	return &job{
		cmd:    cmd,
		status: queued,
	}
}

// start handles running of linux command processes in a go routine
func (j *job) start(id string) error {
	// TODO: set timeout in case process hangs
	// exec.CommandContext() seems to be the suggested way to handle this?

	// channel to catch error inside go routine
	errCh := make(chan error)

	// run process
	go func() {
		// create new Command
		cmd := exec.Command(j.cmd[0], j.cmd[1:]...)

		// set stdout and stderr to write to same buffer
		// mimics (*cmd) CombinedOutput() from os/exec library
		var buf bytes.Buffer
		cmd.Stdout = &buf
		cmd.Stderr = &buf

		// run command
		err := cmd.Start()
		errCh <- err
		if err != nil {
			return
		}

		j.Lock()
		j.status = running
		// store the Pid so stop() can be called later if needed.
		j.pid = cmd.Process.Pid
		j.Unlock()

		cmd.Wait()

		j.Lock()
		defer j.Unlock()
		// store result from process
		j.output = string(buf.Bytes())

		if cmd.ProcessState.Success() {
			j.status = finished
		}
		// BUG: prevents stop() from setting status to canceled when called.
		//  }else {
		// 	j.status = failed
		// }
	}()

	j.Lock()

	// handle error from cmd.Start() in Go routine
	err := <-errCh
	if err != nil {
		j.output = fmt.Errorf("error starting command: %v", err).Error()
		j.status = failed
		return err
	}
	j.Unlock()

	// small delay to let go routine finish writing output for fast executing commands.
	// allows output to be returned on initial response to starting a job
	time.Sleep(50 * time.Millisecond)

	return nil
}

// stop kills process if it is running when called.
func (j *job) stop() (bool, error) {
	j.Lock()
	defer j.Unlock()

	if j.status != running {
		return false, nil
	}

	// lookup by pid and kill when found
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
