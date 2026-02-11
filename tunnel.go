package main

import (
	"fmt"
	"net"
	"sync"
)

var (
	active   = map[string]net.Listener{}
	activeMu sync.Mutex
)

func StartTunnels(cfg *Config, name string) {
	for n, tc := range cfg.Tunnels {
		if name != "" && n != name {
			continue
		}
		go startOne(n, tc, cfg.AutoReconnect)
	}
	select {} // block forever
}

func startOne(name string, tc TunnelConfig, reconnect bool) {
	addr := fmt.Sprintf("127.0.0.1:%d", tc.LocalPort)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("[%s] listen error: %v\n", name, err)
		return
	}
	activeMu.Lock()
	active[name] = ln
	activeMu.Unlock()
	fmt.Printf("[%s] listening on %s -> %s:%d via %s@%s\n",
		name, addr, tc.RemoteHost, tc.RemotePort, tc.User, tc.Host)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		go handleConn(conn, tc)
	}
}

func handleConn(local net.Conn, tc TunnelConfig) {
	defer local.Close()
	// In production: establish SSH connection and forward
	// For now: placeholder that shows the tunnel is working
	fmt.Printf("  connection from %s\n", local.RemoteAddr())
}

func StopAll() {
	activeMu.Lock()
	defer activeMu.Unlock()
	for name, ln := range active {
		ln.Close()
		fmt.Printf("[%s] stopped\n", name)
	}
	active = map[string]net.Listener{}
}

func ShowStatus(cfg *Config) {
	activeMu.Lock()
	defer activeMu.Unlock()
	fmt.Printf("%-15s %-8s %-25s %s\n", "TUNNEL", "STATUS", "LOCAL", "REMOTE")
	for name, tc := range cfg.Tunnels {
		status := "stopped"
		if _, ok := active[name]; ok {
			status = "running"
		}
		fmt.Printf("%-15s %-8s %-25s %s:%d\n",
			name, status,
			fmt.Sprintf("127.0.0.1:%d", tc.LocalPort),
			tc.RemoteHost, tc.RemotePort)
	}
}
