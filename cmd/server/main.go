package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bradyfontenot/ljw/internal/server"
	"github.com/bradyfontenot/ljw/internal/worker"
)

func main() {

	// Removed the auth setup from server.New() to avoid having to error check New() on assignment.

	// init httpserver
	srv := server.New(worker.New())

	// setup Authentication
	if err := srv.SetupTLS(); err != nil {
		fmt.Printf("Problem with authentication setup. Could not start server.\nError: %v\nShutting down...", err)
		os.Exit(1)
	}

	log.Fatal(srv.ListenAndServeTLS("", ""))
}
