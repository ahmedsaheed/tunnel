package main

import (
	"flag"
	"log"
	"os"
)

func ParseFlags() (int, int) {
	from := flag.Int("from", 3000, "Port to forward to")
	to := flag.Int("to", 8080, "Port to listen on")
	flag.Parse()

	if len(os.Args[1:]) == 0 {
		log.Println("No flags provided, using default ports")
		log.Printf("tunnel --from %d --to %d", *from, *to)
	}
	if *from == 0 || *to == 0 {
		log.Println("Please provide the port to forward and the port to listen on")
		return 0, 0
	}

	return *from, *to
}
