package main

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "udp4"
)

func main() {

	serverUDP, err := net.ResolveUDPAddr("udp4", CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP(CONN_TYPE, serverUDP)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	buffer := make([]byte, 1024)

	for {

		n, addr, err := conn.ReadFromUDP(buffer)

		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		fmt.Print("-> ", string(buffer[0:n]))

		data := []byte("hello back!")
		_, err = conn.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
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
