package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// 172.217.24.36
func TraceIt() {

	destination := "www.google.com" // Replace with the target host you want to trace
	maxHops := 30                   // Set the maximum number of hops

	for ttl := 1; ttl <= maxHops; ttl++ {
		//target := fmt.Sprintf("%s:%d", destination, 80)
		target := destination

		conn, err := net.Dial("ip4:icmp", target)
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

	//SendHttpPacket()
	//ExecutePacket()
	//fmt.Println("Directory created successfully")
	//NetworkHops()
	//BasicNetworkConnection()
	TraceIt()
	//TraceRoute()
	//UDPTrace()
}

func UDPTrace() {

	destination := "www.google.com" // Replace with the target host you want to trace
	maxHops := 30                   // Set the maximum number of hops

	for ttl := 1; ttl <= maxHops; ttl++ {
		target := fmt.Sprintf("%s:%d", destination, 33434+ttl)

		conn, err := net.Dial("", target)
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
			fmt.Printf("%2d: * (Timeout) %s \n", ttl, err.Error())
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

func SendHttpPacket() {

	httpRequest := "GET / HTTP/1.1\r\nHost: www.google.com\r\n\r\n"

	// Ethernet layer
	ethernetLayer := layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		DstMAC:       net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
		EthernetType: layers.EthernetTypeIPv4,
	}

	// IP layer
	ipLayer := layers.IPv4{
		Version:    4,
		IHL:        5,
		TOS:        0,
		Length:     20 + uint16(len(httpRequest)),
		Id:         12345,
		Flags:      0,
		FragOffset: 0,
		TTL:        64,
		Protocol:   layers.IPProtocolTCP,
		SrcIP:      net.IP{192, 168, 1, 100},
		DstIP:      net.IP{192, 168, 1, 200},
	}

	// TCP layer
	tcpLayer := layers.TCP{
		SrcPort:  layers.TCPPort(8080),
		DstPort:  layers.TCPPort(80),
		Seq:      12345,
		Ack:      0,
		SYN:      true,
		FIN:      false,
		Window:   65535,
		Checksum: 0,
	}

	// Payload layer (HTTP request)
	payload := gopacket.Payload([]byte(httpRequest))

	// Combine layers into a single packet
	buffer := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{FixLengths: true},
		&ethernetLayer, &ipLayer, &tcpLayer, payload)

	conn, err := net.Dial("tcp", "www.google.com:80")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	_, err = conn.Write(buffer.Bytes())

	if err != nil {
		fmt.Println("Error sending packet:", err)
		return
	}

	conn.Close()

}
