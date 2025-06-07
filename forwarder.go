package main

import (
	"net"
)

type DNSForwarder struct {
	resolverAddr *net.UDPAddr
	sourceIP     string
}

func NewDNSForwarder(resolver string, sourceIP string) (*DNSForwarder, error) {
	addr, err := net.ResolveUDPAddr("udp", resolver)
	if err != nil {
		return nil, err
	}
	return &DNSForwarder{
		resolverAddr: addr,
		sourceIP:     sourceIP,
	}, nil
}

func (f *DNSForwarder) Forward(query []byte) ([]byte, error) {
	var localAddr *net.UDPAddr
	if f.sourceIP != "" {
		var err error
		localAddr, err = net.ResolveUDPAddr("udp", net.JoinHostPort(f.sourceIP, "0"))
		if err != nil {
			return nil, err
		}
	}

	conn, err := net.DialUDP("udp", localAddr, f.resolverAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if _, err := conn.Write(query); err != nil {
		return nil, err
	}

	response := make([]byte, 4096)
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}

	return response[:n], nil
}
