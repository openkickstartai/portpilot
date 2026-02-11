package main

import (
	"net"
	"testing"
	"time"
)

func TestStartSOCKS5(t *testing.T) {
	ln, err := StartSOCKS5("127.0.0.1:0")
	if err != nil { t.Fatal(err) }
	defer ln.Close()
	addr := ln.Addr().String()
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil { t.Fatal("cannot connect to SOCKS proxy") }
	conn.Close()
}
