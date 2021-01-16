package main

import (
	"log"

	"github.com/bradyfontenot/ljw/internal/server"
	"github.com/bradyfontenot/ljw/internal/worker"
)

func main() {

	srv := server.New(worker.New())
	log.Fatal(srv.ListenAndServeTLS("", ""))
}
