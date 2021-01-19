package worker

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

// status values
const (
	queued   = "QUEUED"
	running  = "RUNNING"
	finished = "FINISHED"
	canceled = "CANCELED"
	failed   = "FAILED"
	timeout  = "TIMEOUT"
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
func (j *job) start(id string) {

	go func() {

		// timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// init command
		cmd := exec.CommandContext(ctx, j.cmd[0], j.cmd[1:]...)
		// Create a Process Group ID so call to kill also kills child process
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		// set stdout and stderr to write to single buffer
		var buf bytes.Buffer
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		// run command
		err := cmd.Start()
		if err != nil {
			cancel()
			j.Lock()
			j.output = err.Error()
			j.status = failed
			j.Unlock()
			return
		}

		j.Lock()
		j.status = running
		// store the Pid so stop() can be called later if needed.
		j.pid = cmd.Process.Pid
		j.Unlock()

		// Blocks until process finishes or is killed
		err = cmd.Wait()
		cancel()
		j.Lock()
		defer j.Unlock()

		// record stdout/stderr
		j.output = string(buf.Bytes())
		fmt.Println(err) // debug temp

		// handle premature exit.
		// tbh not sure most reliable way to check exit was due to kill signal.
		// godocs say -1 can mean terminated by signal or that process hasn't exited
		// but if Wait() is no longer blocking then proc should be finished or killed.
		// Also considered using `err.Error() == "signal: killed"`` instead of exit code:
		// not sure if there's a way to grab the signal sent and compare based on Signal type?
		if cmd.ProcessState.ExitCode() == -1 && ctx.Err() != context.DeadlineExceeded {
			j.status = canceled
			return
		} else if ctx.Err() == context.DeadlineExceeded {
			j.status = timeout
			return
		}

		if cmd.ProcessState.Success() {
			j.status = finished
		} else {
			j.status = failed
		}
	}()

	// small delay to let go routine finish writing output for fast executing commands.
	// allows output to be returned on initial response to starting a job
	time.Sleep(50 * time.Millisecond)

}

// stop kills process if it is running when called.
func (j *job) stop() (bool, error) {
	j.RLock()
	defer j.RUnlock()

	if j.status != running {
		return false, nil
	}

	// kill process group
	if err := syscall.Kill(-j.pid, syscall.SIGKILL); err != nil {
		return false, fmt.Errorf("could not kill process. error: %v", err)
	}

	return true, nil
}
