package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync/atomic"
)

var (
	inPackets  uint64
	inBytes    uint64
	outPackets uint64
	outBytes   uint64
)

func main() {
	// Get environment variables
	resolver := os.Getenv("RESOLVER")
	if resolver == "" {
		log.Fatal("Environment variable RESOLVER is not set")
	}

	bindIP := os.Getenv("BIND_IP")
	if bindIP == "" {
		bindIP = "0.0.0.0"
	}

	sourceIP := os.Getenv("SOURCE_IP")
	// SOURCE_IP is optional; default is system-chosen

	// Listen on the specified IP and UDP port 53
	addr := net.JoinHostPort(bindIP, "53")
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to start UDP server: %v", err)
	}
	defer conn.Close()

	log.Printf("DNS resolver running on %s and forwarding to %s (source IP: %s)", addr, resolver, sourceIP)

	buffer := make([]byte, 4096)
	// Start HTTP metrics server
	go startHTTPServer()

	forwarder, err := NewDNSForwarder(resolver, sourceIP)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}
		atomic.AddUint64(&inPackets, 1)
		atomic.AddUint64(&inBytes, uint64(n))

		go handleDNSRequest(buffer[:n], clientAddr, conn, forwarder)
	}
}

func handleDNSRequest(query []byte, clientAddr *net.UDPAddr, serverConn *net.UDPConn, forwarder *DNSForwarder) {
	response, err := forwarder.Forward(query)
	if err != nil {
		log.Printf("DNS forwarding error: %v", err)
		return
	}

	atomic.AddUint64(&outPackets, 1)
	atomic.AddUint64(&outBytes, uint64(len(response)))

	_, err = serverConn.WriteToUDP(response, clientAddr)
	if err != nil {
		log.Printf("Failed to send response: %v", err)
	}
}

func startHTTPServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "DNS Metrics\n")
		fmt.Fprintf(w, "-----------\n")
		fmt.Fprintf(w, "Incoming Packets: %d\n", atomic.LoadUint64(&inPackets))
		fmt.Fprintf(w, "Incoming Bytes:   %d\n", atomic.LoadUint64(&inBytes))
		fmt.Fprintf(w, "Outgoing Packets: %d\n", atomic.LoadUint64(&outPackets))
		fmt.Fprintf(w, "Outgoing Bytes:   %d\n", atomic.LoadUint64(&outBytes))
	})

	log.Println("HTTP metrics server listening on :80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
