package main

import (
	"fmt"
	"os"
	"strings"
)

var version = "0.1.0"

var commands = map[string]string{
	"up":      "Start tunnels (all or by name)",
	"down":    "Stop all active tunnels",
	"status":  "Show tunnel status",
	"version": "Print version",
	"help":    "Show this help message",
}

// ordered for display
var commandOrder = []string{"up", "down", "status", "version", "help"}

func usage() {
	fmt.Fprintf(os.Stderr, "portpilot %s — SSH tunnel manager\n\n", version)
	fmt.Fprintf(os.Stderr, "Usage:\n  portpilot [-c config] <command> [args]\n\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	for _, cmd := range commandOrder {
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", cmd, commands[cmd])
	}
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fmt.Fprintf(os.Stderr, "  -c, --config <path>   Config file (default: ~/.portpilot.yml)\n")
	fmt.Fprintf(os.Stderr, "  -h, --help            Show help\n")
	fmt.Fprintf(os.Stderr, "  -v, --version         Print version\n")
	fmt.Fprintf(os.Stderr, "\nEnvironment:\n")
	fmt.Fprintf(os.Stderr, "  PORTPILOT_CONFIG      Override default config path\n")
}

func main() {
	args := os.Args[1:]
	cfgOverride := ""

	// Parse global flags before subcommand
	var filtered []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-c", "--config":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "error: -c/--config requires a path argument")
				os.Exit(1)
			}
			i++
			cfgOverride = args[i]
		case "-h", "--help":
			usage()
			os.Exit(0)
		case "-v", "--version":
			fmt.Printf("portpilot %s\n", version)
			os.Exit(0)
		default:
			if strings.HasPrefix(args[i], "-") {
				fmt.Fprintf(os.Stderr, "error: unknown flag %q\n\nRun 'portpilot --help' for usage.\n", args[i])
				os.Exit(1)
			}
			filtered = append(filtered, args[i])
		}
	}

	if len(filtered) == 0 {
		usage()
		os.Exit(1)
	}

	cmd := filtered[0]
	if _, ok := commands[cmd]; !ok {
		fmt.Fprintf(os.Stderr, "error: unknown command %q\n", cmd)
		fmt.Fprintf(os.Stderr, "valid: %s\n\nRun 'portpilot help' for usage.\n", strings.Join(commandOrder, ", "))
		os.Exit(1)
	}

	switch cmd {
	case "version":
		fmt.Printf("portpilot %s\n", version)
		return
	case "help":
		usage()
		return
	}

	cfgPath := cfgOverride
	if cfgPath == "" {
		cfgPath = configPath()
	}

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		fmt.Fprintf(os.Stderr, "hint: create %s or specify with -c <path>\n", cfgPath)
		os.Exit(1)
	}

	switch cmd {
	case "up":
		name := ""
		if len(filtered) > 1 {
			name = filtered[1]
		}
		StartTunnels(cfg, name)
	case "down":
		StopAll()
	case "status":
		ShowStatus(cfg)
	}
}

func configPath() string {
	if p := os.Getenv("PORTPILOT_CONFIG"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: cannot determine home dir: %v, falling back to .portpilot.yml\n", err)
		return ".portpilot.yml"
	}
	return home + "/.portpilot.yml"
}
