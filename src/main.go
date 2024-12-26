package main

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	from, to := ParseFlags()
	if from == 0 || to == 0 {
		os.Exit(1)
	}
	serverPort := strconv.Itoa(to)
	targetPort := strconv.Itoa(from)
	listener, err := net.Listen("tcp", "localhost:"+serverPort)

	if err != nil {
		log.Fatalf("Error starting listener: %v", err)
	}
	defer listener.Close()
	log.Println("Tunnel created on localhost:" + serverPort)
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(clientConn, targetPort)
	}
}

func handleConnection(clientConn net.Conn, targetPort string) {
	defer clientConn.Close()

	// Connect to the target server on port 4000
	targetConn, err := net.Dial("tcp", "localhost:"+targetPort)
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
