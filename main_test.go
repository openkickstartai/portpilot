package main

import (
	"os"
	"testing"
)

func TestConfigPathDefault(t *testing.T) {
	os.Unsetenv("PORTPILOT_CONFIG")
	p := configPath()
	if p == "" {
		t.Fatal("configPath() returned empty string")
	}
	if p == "/.portpilot.yml" {
		t.Fatal("configPath() returned root-relative path, UserHomeDir likely failed silently")
	}
}

func TestConfigPathEnvOverride(t *testing.T) {
	os.Setenv("PORTPILOT_CONFIG", "/tmp/test-portpilot.yml")
	defer os.Unsetenv("PORTPILOT_CONFIG")
	p := configPath()
	if p != "/tmp/test-portpilot.yml" {
		t.Fatalf("expected env override, got %q", p)
	}
}

func TestCommandsMapConsistency(t *testing.T) {
	// every entry in commandOrder must exist in commands map
	for _, cmd := range commandOrder {
		if _, ok := commands[cmd]; !ok {
			t.Errorf("commandOrder has %q but commands map does not", cmd)
		}
	}
	// every entry in commands map must be in commandOrder
	for cmd := range commands {
		found := false
		for _, c := range commandOrder {
			if c == cmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("commands map has %q but commandOrder does not", cmd)
		}
	}
}
