package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func ExecutePacket() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: go-traceroute <target>")
		os.Exit(1)
	}

	target := os.Args[1]
	maxHops := 30
	var ttl uint8 = 1

	for i := 0; i < maxHops; i++ {

		// Create an ICMP echo request packet
		icmp := &layers.ICMPv4{
			TypeCode: layers.ICMPv4TypeEchoRequest,
			Id:       12345,
			Seq:      1,
		}

		payload := gopacket.Payload([]byte("Hello, world!"))
		buffer := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}

		// Set the TTL for the IP header
		ip := &layers.IPv4{
			TTL: ttl,
		}
		ip.SerializeTo(buffer, opts)
		icmp.SerializeTo(buffer, opts)
		payload.SerializeTo(buffer, opts)

		// Send the packet
		conn, err := net.Dial("tcp", target)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer conn.Close()

		_, err = conn.Write(buffer.Bytes())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Receive the response
		conn.SetReadDeadline(time.Now().Add(time.Second))
		data := make([]byte, 1500)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Printf("%d hops away: ???\n", ttl)
			ttl++
			continue
		}

		// 		// Decode a packet
		// packet := gopacket.NewPacket(myPacketData, layers.LayerTypeEthernet, gopacket.Default)
		// // Get the TCP layer from this packet
		// if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		//   fmt.Println("This is a TCP packet!")
		//   // Get actual TCP data from this layer
		//   tcp, _ := tcpLayer.(*layers.TCP)
		//   fmt.Printf("From src port %d to dst port %d\n", tcp.SrcPort, tcp.DstPort)
		// }
		// // Iterate over all layers, printing out each layer type
		// for _, layer := range packet.Layers() {
		//   fmt.Println("PACKET LAYER:", layer.LayerType())
		// }

		// Parse the response
		packet := gopacket.NewPacket(data[:n], layers.LayerTypeIPv4, gopacket.Default)
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			fmt.Println("Invalid IP packet")
			os.Exit(1)
		}

		ip = ipLayer.(*layers.IPv4)
		srcIP := ip.SrcIP.String()

		icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
		if icmpLayer == nil {
			fmt.Println("Invalid ICMP packet")
			os.Exit(1)
		}

		icmp = icmpLayer.(*layers.ICMPv4)
		if icmp.TypeCode.Type() == layers.ICMPv4TypeTimeExceeded {
			fmt.Printf("%d you're hops away: %s\n", ttl, srcIP)
			ttl++
		} else if icmp.TypeCode.Type() == layers.ICMPv4TypeEchoReply {
			fmt.Printf("%d hops away: %s\n", ttl, srcIP)
			break
		} else {
			fmt.Println("Unknown ICMP packet type")
			os.Exit(1)
		}
	}

}
