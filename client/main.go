package main

import (
	"flag"
	"io"
	"log"
	"net"
)

func startTunnel(localPort string, serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	localListener, err := net.Listen("tcp", localPort)
	if err != nil {
		log.Fatalf("Error starting local server: %v", err)
	}
	defer localListener.Close()

	for {
		clientConn, err := localListener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		log.Printf("Accepted connection from %v", clientConn.RemoteAddr())

		go func() {
			defer clientConn.Close()
			io.Copy(conn, clientConn)
		}()
	}
}

func main() {
	localPort := flag.String("local", ":3030", "Localhost port to expose")
	serverAddr := flag.String("server", "localhost:8080", "Public server address")
	flag.Parse()

	startTunnel(*localPort, *serverAddr)
}
