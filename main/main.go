package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {

	// tunnel --from <target_port> --to <server_port>
	from := flag.Int("from", 3000, "Port to forward to")
	to := flag.Int("to", 8080, "Port to listen on")
	flag.Parse()
	fmt.Println(flag.Args())

	if len(os.Args[1:]) == 0 {
		log.Println("No flags provided, using default ports")
		log.Printf("tunnel --from %d --to %d", *from, *to)
	}
	if *from == 0 || *to == 0 {
		log.Println("Please provide the port to forward and the port to listen on")
		log.Fatalf("Usage: tunnel --from <target_port> --to <server_port>")
	}

	serverPort := strconv.Itoa(*to)
	targetPort := strconv.Itoa(*from)
	listener, err := net.Listen("tcp", "localhost:"+serverPort)

	if err != nil {
		log.Fatalf("Error starting listener: %v", err)
	}
	defer listener.Close()

	log.Println("Tunnel is running: localhost:" + serverPort + "-> localhost:" + targetPort)
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
