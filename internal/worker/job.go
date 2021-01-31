package worker

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
	"syscall"
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

	cmd := exec.Command(j.cmd[0], j.cmd[1:]...)
	// Create a Process Group ID so call to terminate also kills child process
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmdReader, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	done := make(chan bool)
	sc := make(chan bool)
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		sc <- true
		for scanner.Scan() {
			line := scanner.Text()
			j.Lock()
			j.output = append(j.output, line+"\n")
			j.Unlock()
		}
		done <- true
	}()

	<-sc
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

}

// stop kills process if it is running when called and returns true.
// otherwise simply returns false.
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
