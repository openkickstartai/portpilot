package main

import (
	"fmt"
	"os"
)

var version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: portpilot <up|down|status|version>")
		os.Exit(1)
	}
	cfg, err := LoadConfig(configPath())
	if err != nil && os.Args[1] != "version" {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "up":
		name := ""
		if len(os.Args) > 2 { name = os.Args[2] }
		StartTunnels(cfg, name)
	case "down":
		StopAll()
	case "status":
		ShowStatus(cfg)
	case "version":
		fmt.Printf("portpilot %s\n", version)
	default:
		fmt.Fprintf(os.Stderr, "unknown: %s\n", os.Args[1])
	}
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return home + "/.portpilot.yml"
}
