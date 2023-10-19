package main

import (
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "3333"
	TYPE = "ip4"
)

func main() {

	ipServer, err := net.ResolveIPAddr(TYPE, HOST)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialIP("ip4:icmp", nil, ipServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("This is a message"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	// getting data from server
	// received := make([]byte, 1024)
	// _, err = conn.Read(received)
	// if err != nil {
	// 	println("Read data failed:", err.Error())
	// 	os.Exit(1)
	// }

	// println("Received message:", string(received))

	conn.Close()
}
