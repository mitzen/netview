package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	destination := "localhost:3333"
	target := destination
	conn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Printf("%2d: * (Timeout)\n")
	}
	defer conn.Close()

	// Set a timeout for reading responses
	conn.SetDeadline(time.Now().Add(3 * time.Second))

	startTime := time.Now()
	_, err = conn.Write([]byte("Hello, World!")) // Sending a UDP packet
	if err != nil {
		fmt.Printf("%2d: * (Timeout)\n")
	}

	reply := make([]byte, 1024)
	_, err = conn.Read(reply) // Receive a response
	elapsedTime := time.Since(startTime)

	if err != nil {
		fmt.Printf("%2d: * (Timeout)\n")
	} else {
		remoteAddr := conn.RemoteAddr()
		fmt.Printf("%2d: %s %.2f ms\n", remoteAddr, float64(elapsedTime.Milliseconds()))
	}

	if destination == conn.RemoteAddr().String() {
		fmt.Printf("Trace completed in %d hops.\n")
	}
}
