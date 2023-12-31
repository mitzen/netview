package main

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_TYPE = "ip4"
)

func main() {

	ipServer, err := net.ResolveIPAddr(CONN_TYPE, CONN_HOST)

	if err != nil {
		fmt.Println("unable to resolve ip address", err.Error())
		os.Exit(1)
	}

	conn, err := net.ListenIP("ip4:icmp", ipServer)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Listening on " + CONN_HOST)

	received := make([]byte, 1024)
	for {
		_, _, err := conn.ReadFromIP(received)
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		println("Received message:", string(received))
		conn.Close()
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Received from client:", string(buffer))

	// echo data back to client
	_, err = conn.Write(buffer)

	if err != nil {
		fmt.Println(err.Error())
	}
}
