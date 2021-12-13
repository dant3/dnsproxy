package main

import (
	"context"
	"fmt"
	"net"
)

type DnsProxyServer struct {
	ctx    context.Context
	port   int
	config *Config
}

func NewDnsProxyServer(port int, config *Config) DnsProxyServer {
	return DnsProxyServer{
		ctx:    context.Background(),
		port:   port,
		config: config,
	}
}

const maxBufferSize = 512

func (server DnsProxyServer) run() {
	addr := new(net.UDPAddr)
	addr.Port = server.port

	conn, err := net.ListenUDP("udp", addr)
	exitOnError(err, "Failed to open socket: %v")

	defer conn.Close()

	buffer := make([]byte, maxBufferSize)

	fmt.Println("DNS server is running on port", server.port)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		exitOnError(err, "Failed to read from socket: %v")

		fmt.Printf("packet-received: bytes=%d from=%s\n", n, addr.String())
		response, err := process(buffer, conn, addr, server.config)
		if response != nil {
			conn.WriteTo(response, addr)
		}
	}
}

func process(packet []byte, conn net.PacketConn, remoteAddr net.Addr, config *Config) ([]byte, error) {
	dnsRequest := DecodeRequest(packet)

	if config.isBlacklisted(dnsRequest.question.qname) {
		fmt.Println("Blacklisted address:", dnsRequest.question.qname)
		response := rejectResponse(dnsRequest)
		return EncodePacket(response), nil
	} else {
		fmt.Println("Whitelisted address:", dnsRequest.question.qname)
		return proxyTo(packet, config.nameserver)
	}
}

func rejectResponse(request DnsRequest) DnsPacket {
	header := request.header
	header.qr = true
	header.ancount = 0
	header.rcode = uint8(0b0101)

	questions := make([]DnsQuestion, 1)
	questions[0] = request.question
	response := DnsPacket{
		header:    header,
		questions: questions,
	}
	return response
}

func proxyTo(packet []byte, relayAddress string) (response []byte, err error) {
	conn, err := net.Dial("udp", relayAddress+":53")
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	_, err = conn.Write(packet)
	_, err = conn.Read(packet)

	if err != nil {
		return nil, err
	} else {
		return packet, nil
	}
}
