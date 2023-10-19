package main

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
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
