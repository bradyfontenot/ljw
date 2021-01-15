package main

import (
	"fmt"
	"os"

	"github.com/bradyfontenot/ljw/internal/client"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Print("No arguments supplied. Must supply at least one argument\n")
		return
	}
	// collect command line args
	appCommand := os.Args[1]
	args := os.Args[2:]

	// create new client
	c, err := client.New()
	if err != nil {
		fmt.Printf("Problem Authenticating.\nError: %s\nShutting down...\n", err)
		os.Exit(1)
	}

	switch appCommand {
	case "list":
		list(c)
	case "status":
		status(c, args[0])
	case "start":
		start(c, args[0:])
	case "stop":
		stop(c, args[0])
	case "log":
		log(c, args[0])
	}
}

func list(c *client.Client) {
	err := c.ListRunningJobs()
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}

}

func status(c *client.Client, id string) {
	err := c.JobStatus(id)
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}
}

func start(c *client.Client, args []string) {
	err := c.StartJob(args)
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}
}

func stop(c *client.Client, id string) {
	err := c.StopJob(id)
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}
}

func log(c *client.Client, id string) {
	err := c.GetJobLog(id)
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}
}

//########################## RAW FUNCTIONS #########################################

// func main() {

// 	c, err := client.New()
// 	if err != nil {
// 		fmt.Printf("Problem Authenticating.\nError: %s\nShutting down...\n", err)
// 		os.Exit(1)
// 	}

// 	fmt.Print("Client started\n")
// 	cmd := make([]string, 0)
// 	cmd = append(cmd, "sleep")
// 	cmd = append(cmd, "3")

// 	// Temporary Calls to check request/response data
// 	// will be replaced with a cli
// 	for i := 0; i < 5; i++ {
// 		c.StartJob(cmd)
// 		time.Sleep(time.Duration(1) * time.Second)
// 	}
// c.ListRunningJobs()
// 	c.JobStatus("3")
// 	c.StopJob("3")
// 	c.JobStatus("3")
// 	c.GetJobLog("3")
// }
