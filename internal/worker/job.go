package worker

import (
	"bytes"
	"fmt"
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

	// out, err := exec.Command("grep", "main").Output()
	c := exec.Command("grep", "brady", "go.mod") //, "brady go.mod")
	var out bytes.Buffer
	c.Stdout = &out
	err := c.Run()
	if err != nil {
		j.status = failed
		fmt.Println(err)
	}
	fmt.Printf("Results: %q\n", out.String())

	return nil
}

func (j *Job) stop() error {

	// kill proc

	j.status = canceled

	return nil
}
