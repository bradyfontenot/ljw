package main

import (
	"fmt"
	"os"

	"github.com/bradyfontenot/ljw/internal/client"
)

func main() {

	c, err := client.New()
	if err != nil {
		fmt.Printf("Problem with authentication setup. Could not start client.\nError: %v\nShutting down...", err)
		os.Exit(1)
	}

	fmt.Print("Client started\n")

	// Temporary Calls to check request/response data
	// will be replaced with a cli
	for i := 0; i < 5; i++ {
		c.StartJob()
	}
	c.ListRunningJobs()
	c.JobStatus("3")
	c.StopJob("3")
	c.JobStatus("3")
	c.GetJobLog("3")
}
