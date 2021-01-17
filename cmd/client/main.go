package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/bradyfontenot/ljw/internal/client"
)

func main() {

	// check if args are supplied
	if len(os.Args) <= 1 {
		fmt.Println("\nNo arguments supplied. Must supply at least one argument\n")
		printUsage()
		return
	}
	// collect command line args
	appCommand := os.Args[1]
	args := os.Args[2:]

	// start client and attempt to connect to server
	c, err := client.New()
	if err != nil {
		fmt.Printf("Problem Authenticating.\nError: %s\nShutting down...\n", err)
		return
	}

	switch appCommand {
	case "list":
		list(c, args[0:])
	case "status":
		status(c, args[0:])
	case "start":
		start(c, args[0:])
	case "stop":
		stop(c, args[0:])
	case "log":
		log(c, args[0:])
	default:
		printUsage()
	}
}

func list(c *client.Client, args []string) {
	// validate args exist
	if len(args) > 0 {
		fmt.Println("\nToo many args. list takes no arguments.\n")
		printUsage()
		return
	}

	err := c.ListRunningJobs()
	if err != nil {
		printError(err)
		return
	}
}

func status(c *client.Client, args []string) {
	// validate only one arg supplied for id.
	id, err := processID(args)
	if err != nil {
		return
	}

	err = c.JobStatus(id)
	if err != nil {
		printError(err)
		return
	}
}

func start(c *client.Client, args []string) {
	// validate args exist
	if len(args) < 1 {
		fmt.Println("\nNo linux command supplied. Must supply a command\n")
		printUsage()
		return
	}

	err := c.StartJob(args)
	if err != nil {
		printError(err)
		return
	}
}

func stop(c *client.Client, args []string) {
	// validate only one arg supplied for id
	id, err := processID(args)
	if err != nil {
		return
	}

	err = c.StopJob(id)
	if err != nil {
		printError(err)
		return
	}
}

func log(c *client.Client, args []string) {
	// validate only one arg supplied for id.
	id, err := processID(args)
	if err != nil {
		return
	}

	err = c.GetJobLog(id)
	if err != nil {
		printError(err)
		return
	}
}

// printUsage prints usage instructions
func printUsage() {
	fmt.Println("[USAGE]")
	fmt.Printf(" list\n start \t<linux cmd>\n status\t<job id>\n stop \t<job id>\n log \t<job id>\n\n")
}

// processID ensures only one arg was supplied for id
func processID(args []string) (string, error) {
	if len(args) < 1 {
		fmt.Println("\nNo id supplied.\n")
		printUsage()
		return "", errors.New("no argument supplied")
	} else if len(args) > 1 {
		fmt.Println("\nToo many args or id has spaces. Please Supply only one id at a time.\n")
		printUsage()
		return "", errors.New("too many arguments supplied. only needs 1")
	}

	return args[0], nil
}

// printError formats and prints errors to screen
func printError(err error) {
	fmt.Printf("\n[Error]\n%v \n\n", err)
}
