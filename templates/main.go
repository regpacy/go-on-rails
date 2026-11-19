package main

import (
	"log"
)

func main() {
	server := NewServer()

	if err := server.start(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
