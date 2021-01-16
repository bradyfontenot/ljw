package worker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

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

func (j *job) Status() string {
	j.RLock()
	defer j.RUnlock()
	return j.status
}

func (j *job) Cmd() []string {
	j.RLock()
	defer j.RUnlock()
	return j.cmd
}

func (j *job) Output() string {
	j.RLock()
	defer j.RUnlock()
	return j.output
}

// start handles running of linux command processes in a go routine
func (j *job) start(id string) error {
	// set timeout in case process hangs
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// init new Command
	cmd := exec.Command(j.cmd[0], j.cmd[1:]...)
	// cmd := exec.CommandContext(ctx, j.cmd[0], j.cmd[1:]...)
	// set stdout and stderr to write to same buffer
	// mimics (*cmd) CombinedOutput() from os/exec library
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	// channel to catch error inside go routine
	errCh := make(chan error)
	// run process
	go func() {
		// run command
		err := cmd.Start()
		errCh <- err
		if err != nil {
			return
		}

		j.Lock()
		// set job status
		j.status = running
		// store the Pid so stop() can be called later if needed.
		j.pid = cmd.Process.Pid
		j.Unlock()

		// wait for process to end before accessing buffer below
		cmd.Wait()

		j.Lock()
		// store result from process
		j.output = string(buf.Bytes())

		if cmd.ProcessState.Success() {
			j.status = finished
		}
		//  else {
		// 	j.status = failed
		// }

		j.Unlock()

	}()

	// handle error from cmd.Start() in Go routine
	err := <-errCh
	if err != nil {
		j.Lock()
		j.output = fmt.Errorf("error starting command: %v", err).Error()
		j.status = failed
		j.Unlock()
		return err
	}
	// if ctx.Err() == context.DeadlineExceeded {
	// 	j.status = failed
	// 	j.output = ctx.Err().Error()
	// 	fmt.Print("ctx died")
	// 	return ctx.Err()
	// }

	// small delay to let go routin finish writing output for fast executing commands.
	// allows output to be returned on initial response to starting a job
	time.Sleep(300 * time.Millisecond)

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
