package main

import (
	"fmt"

	"github.com/bradyfontenot/ljw/internal/client"
)

func main() {

	fmt.Print("Client Starting...\n")

	c := client.New()

	fmt.Print("Client started\n")

	// Temporary Calls to check request/response data
	for i := 0; i < 5; i++ {
		c.StartJob()
	}
	c.ListJobs()
	c.JobStatus("3")
	c.StopJob("3")
	c.JobStatus("3")
	c.GetJobLog("3")
}
