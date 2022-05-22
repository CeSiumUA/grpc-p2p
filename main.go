package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]

	if args[0] == "c" {
		log.Println("starting client...")
		StartP2P()
	} else if args[0] == "s" {
		log.Println("starting server...")
		Serve()
	}
	log.Println("node work finished")
}
