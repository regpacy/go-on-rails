package main

import (
	"log"
)

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}
	defer server.deinit()

	if err := server.start(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
