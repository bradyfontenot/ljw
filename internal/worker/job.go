package worker

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
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
)

type job struct {
	cmd    []string
	status string
	output []string
	pid    int
	sync.RWMutex
}

func newJob(cmd []string) *job {
	return &job{
		cmd:    cmd,
		status: queued,
	}
}

// start handles running of linux command processes
func (j *job) start() {
	cmdString := strings.Join(j.cmd[:], " ")
	cmd := exec.Command("sh", "-c", cmdString)
	// Create a Process Group ID so call to terminate also kills child process
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmdReader, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	done := make(chan int)
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			j.Lock()
			j.output = append(j.output, line+"\n")
			j.Unlock()
		}
		done <- 1
	}()

	err = cmd.Start()
	if err != nil {
		j.Lock()
		j.output = append(j.output, err.Error())
		j.status = failed
		j.Unlock()
		return
	}

	j.Lock()
	j.status = running
	// store the Pid so stop() can be called later if needed.
	j.pid = cmd.Process.Pid
	j.Unlock()
	go func() {
		<-done
		err = cmd.Wait()
		j.Lock()
		defer j.Unlock()
		if err != nil {
			j.output = append(j.output, "Error: "+err.Error())
		}

		// handle premature exit.
		// tbh not sure most reliable way to check exit was due to kill signal.
		// godocs say -1 can mean terminated by signal or that process hasn't exited
		// but if Wait() is no longer blocking then proc should be finished or killed.
		// Also considered using `err.Error() == "signal: killed"`` instead of exit code:
		// not sure if there's a way to grab the signal sent and compare based on Signal type?
		if cmd.ProcessState.ExitCode() == -1 {
			j.status = canceled
			return
		}
		if cmd.ProcessState.Success() {
			j.status = finished
		} else {
			j.status = failed
		}
	}()
	time.Sleep(15 * time.Millisecond)
}

// stop kills process if it is running when called.
func (j *job) stop() (bool, error) {
	j.RLock()
	defer j.RUnlock()

	if j.status != running {
		return false, nil
	}

	if err := syscall.Kill(-j.pid, syscall.SIGTERM); err != nil {
		return false, fmt.Errorf("could not terminate process. error: %v", err)
	}

	return true, nil
}

func (j *job) Cmd() []string {
	j.RLock()
	defer j.RUnlock()
	return j.cmd
}

func (j *job) Status() string {
	j.RLock()
	defer j.RUnlock()
	return j.status
}

func (j *job) Output() []string {
	j.RLock()
	defer j.RUnlock()
	return j.output
}
