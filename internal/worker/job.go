package worker

import (
	"fmt"
	"log"
	"os/exec"
)

const (
	running   = "RUNNING"
	queued    = "QUEUED"
	completed = "COMPLETED"
	canceled  = "CANCELED"
	failed    = "FAILED"
)

// Job is one linux job
type Job struct {
	cmd    string
	status string
	output string
}

// New creates a new Job
func newJob(cmd string) *Job {
	return &Job{
		cmd:    cmd,
		status: running,
		// output: "none",
	}
}

// TODO:
//	* handle invalid cmd(don't kill server. just set status to Failed)
//	* Run in go routine ya?

func (j *Job) start() error {
	// set job status
	j.status = running

	stdoutStderr, err := exec.Command("ls").CombinedOutput()
	if err != nil {
		log.Printf("%v\n", stdoutStderr)
	}
	fmt.Printf("%s", stdoutStderr)

	j.output = string(stdoutStderr)
	j.status = completed

	return nil
}

func (j *Job) stop() error {

	// kill proc

	j.status = canceled

	return nil
}
