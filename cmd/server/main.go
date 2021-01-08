package main

import (
	"log"
	"net/http"

	"github.com/bradyfontenot/ljw/internal/server"
)

var (
	port = ":8080"
)

func main() {

	srv := server.New(port)
	if err := http.ListenAndServe(port, srv); err != nil {
		log.Fatalf("Could not connect to server on port %v. \n %v", port, err)
	}
}
