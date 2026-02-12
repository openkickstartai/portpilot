package main

import (
	"fmt"
	"io"
	"net"
)

// StartSOCKS5 starts a local SOCKS5 proxy that tunnels through SSH.
func StartSOCKS5(listenAddr string) (net.Listener, error) {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("socks listen: %w", err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil { return }
			go handleSOCKS(conn)
		}
	}()
	return ln, nil
}

func handleSOCKS(conn net.Conn) {
	defer conn.Close()
	// Read SOCKS5 greeting
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil || n < 3 || buf[0] != 0x05 {
		return
	}
	// No auth response
	conn.Write([]byte{0x05, 0x00})
	// Read connect request
	n, err = conn.Read(buf)
	if err != nil || n < 7 { return }
	var host string
	var port int
	switch buf[3] {
	case 0x01: // IPv4
		host = fmt.Sprintf("%d.%d.%d.%d", buf[4], buf[5], buf[6], buf[7])
		port = int(buf[8])<<8 | int(buf[9])
	case 0x03: // Domain
		l := int(buf[4])
		host = string(buf[5:5+l])
		port = int(buf[5+l])<<8 | int(buf[6+l])
	default:
		return
	}
	// Connect to target (in production: through SSH)
	target, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		conn.Write([]byte{0x05, 0x05, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		fmt.Printf("Failed to connect to %s:%d - %v\n", host, port, err)
		return
	}
	defer target.Close()
	conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	go io.Copy(target, conn)
	io.Copy(conn, target)
}