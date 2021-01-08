package main

import (
	"fmt"

	"github.com/bradyfontenot/ljw/internal/client"
)

func main() {

	fmt.Print("Client Starting...\n")

	c := client.New()
	fmt.Print("Client started\n")
	c.TestGet()
}
