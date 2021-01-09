package main

import (
	"log"

	"github.com/bradyfontenot/ljw/internal/server"
)

func main() {

	srv := server.New()
	log.Fatal(srv.ListenAndServeTLS("", ""))
}
