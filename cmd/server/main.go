package main

import (
	"log"

	"github.com/bradyfontenot/ljw/internal/server"
	"github.com/bradyfontenot/ljw/internal/worker"
)

func main() {

	wkr := worker.New()
	srv := server.New(wkr)
	log.Fatal(srv.ListenAndServeTLS("", ""))
}
