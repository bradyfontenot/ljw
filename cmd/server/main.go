package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bradyfontenot/ljw/internal/server"
	"github.com/bradyfontenot/ljw/internal/worker"
)

func main() {

	srv, err := server.New(worker.New())
	if err != nil {
		fmt.Printf("Problem with authentication setup. Could not start server.\nError: %v\nShutting down...", err)
		os.Exit(1)
	}

	log.Fatal(srv.ListenAndServeTLS("", ""))
}
