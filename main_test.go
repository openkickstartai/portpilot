package main

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestValidCommands_AllExpected(t *testing.T) {
	expected := []string{"up", "down", "status", "version"}
	for _, cmd := range expected {
		if !validCommands[cmd] {
			t.Errorf("expected %q to be a valid command", cmd)
		}
	}
}

func TestValidCommands_RejectsUnknown(t *testing.T) {
	invalid := []string{"", "start", "stop", "restart", "UP", "Up", " up", "up ", "--help"}
	for _, cmd := range invalid {
		if validCommands[cmd] {
			t.Errorf("expected %q to NOT be a valid command", cmd)
		}
	}
}

func TestValidCommands_Count(t *testing.T) {
	if len(validCommands) != 4 {
		t.Errorf("expected exactly 4 valid commands, got %d", len(validCommands))
	}
}

func TestConfigPath_ContainsFilename(t *testing.T) {
	p := configPath()
	if !strings.HasSuffix(p, ".portpilot.yml") {
		t.Errorf("configPath() = %q, want suffix .portpilot.yml", p)
	}
}

func TestConfigPath_IsAbsolute(t *testing.T) {
	p := configPath()
	// On all platforms the home dir should yield an absolute path
	if runtime.GOOS == "windows" {
		if len(p) < 3 || p[1] != ':' {
			t.Errorf("expected absolute path on windows, got %q", p)
		}
	} else {
		if !strings.HasPrefix(p, "/") {
			t.Errorf("expected absolute path, got %q", p)
		}
	}
}

func TestConfigPath_MatchesHomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}
	expected := home + "/.portpilot.yml"
	got := configPath()
	if got != expected {
		t.Errorf("configPath() = %q, want %q", got, expected)
	}
}

func TestVersion_NotEmpty(t *testing.T) {
	if version == "" {
		t.Error("version should not be empty")
	}
}

func TestVersion_Format(t *testing.T) {
	// Semantic versioning: should have at least two dots
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		t.Errorf("version %q does not look like semver", version)
	}
}
