package main

import (
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error starting listener: %v", err)
	}
	defer listener.Close()

	log.Println("Tunnel is running: localhost:8080 -> localhost:3000")

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(clientConn)
	}
}

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	// Connect to the target server on port 4000
	targetConn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Printf("Error connecting to target: %v", err)
		return
	}

	log.Printf("Accepted connection from %v", clientConn.RemoteAddr())

	defer targetConn.Close()
	log.Printf("Connected to target: %v", targetConn.RemoteAddr())

	// Forward data between client and target
	go io.Copy(targetConn, clientConn) // Client to Target
	io.Copy(clientConn, targetConn)    // Target to Client
}
