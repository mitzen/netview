package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func TraceIt() {

	destination := "www.google.com" // Replace with the target host you want to trace
	maxHops := 30                   // Set the maximum number of hops

	for ttl := 1; ttl <= maxHops; ttl++ {
		target := fmt.Sprintf("%s:%d", destination, 33434+ttl)

		conn, err := net.Dial("udp", target)
		if err != nil {
			fmt.Printf("%2d: * (Timeout)\n", ttl)
			continue
		}
		defer conn.Close()

		// Set a timeout for reading responses
		conn.SetDeadline(time.Now().Add(3 * time.Second))

		startTime := time.Now()
		_, err = conn.Write([]byte("Hello, World!")) // Sending a UDP packet
		if err != nil {
			fmt.Printf("%2d: * (Timeout)\n", ttl)
			continue
		}

		reply := make([]byte, 1024)
		_, err = conn.Read(reply) // Receive a response
		elapsedTime := time.Since(startTime)

		if err != nil {
			fmt.Printf("%2d: * (Timeout)\n", ttl)
		} else {
			remoteAddr := conn.RemoteAddr()
			fmt.Printf("%2d: %s %.2f ms\n", ttl, remoteAddr, float64(elapsedTime.Milliseconds()))
		}

		if destination == conn.RemoteAddr().String() {
			fmt.Printf("Trace completed in %d hops.\n", ttl)
			break
		}
	}
}
func main() {

	// address, err := net.ResolveTCPAddr("", ip)
	// if err != nil {
	// 	fmt.Printf("resolve error:", err)
	// }

	ExecutePacket()
	//fmt.Println("Directory created successfully")
	//NetworkHops()
	//BasicNetworkConnection()
	//TraceIt()
}

func BasicNetworkConnection() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run trace.go <destination IP>")
		os.Exit(1)
	}

	target := os.Args[1]

	netconn, err := net.Dial("tcp", target)

	if err != nil {
		fmt.Printf("error dial", err)
	}

	fmt.Printf("closing connections")
	defer netconn.Close()
}

func NetworkHops() {

	target := "google.com"
	maxHops := 30
	ttl := 1

	for i := 0; i < maxHops; i++ {

		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:80", target), time.Second)

		if err == nil {
			fmt.Printf("%d hops away: %s\n", ttl, conn.RemoteAddr().String())
			conn.Close()
			break
		}

		fmt.Printf("%d hops away: ???\n", ttl)
		ttl++
	}
}
