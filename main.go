package main

import (
	"fmt"
	"os"
)

var version = "0.1.0"

// Command lookup map for better performance
var validCommands = map[string]bool{
	"up":      true,
	"down":    true,
	"status":  true,
	"version": true,
}

func main() {
	argc := len(os.Args)
	if argc < 2 {
		fmt.Println("usage: portpilot <up|down|status|version>")
		os.Exit(1)
	}
	
	cmd := os.Args[1]
	if !validCommands[cmd] {
		fmt.Fprintf(os.Stderr, "unknown: %s\n", cmd)
		os.Exit(1)
	}
	
	cfg, err := LoadConfig(configPath())
	if err != nil && cmd != "version" {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}
	
	switch cmd {
	case "up":
		name := ""
		if argc > 2 {
			name = os.Args[2]
		}
		StartTunnels(cfg, name)
	case "down":
		StopAll()
	case "status":
		ShowStatus(cfg)
	case "version":
		fmt.Printf("portpilot %s\n", version)
	}
}

func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: cannot determine home dir: %v\n", err)
		return ".portpilot.yml"
	}
	return home + "/.portpilot.yml"
}
