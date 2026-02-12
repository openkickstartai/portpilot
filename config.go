package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type TunnelConfig struct {
	Host       string `yaml:"host"`
	User       string `yaml:"user"`
	LocalPort  int    `yaml:"local_port"`
	RemoteHost string `yaml:"remote_host"`
	RemotePort int    `yaml:"remote_port"`
	KeyFile    string `yaml:"key"`
}

type Config struct {
	Tunnels        map[string]TunnelConfig `yaml:"tunnels"`
	AutoReconnect  bool                    `yaml:"auto_reconnect"`
	ReconnectDelay string                  `yaml:"reconnect_delay"`
}

func LoadConfig(path string) (*Config, error) {
	log.Printf("Loading configuration from: %s", path)
	
	// Validate file path to prevent directory traversal
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("invalid config path: directory traversal detected")
	}
	
	// Check if file exists and is readable
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", cleanPath)
	}
	
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		log.Printf("Failed to read config file %s: %v", cleanPath, err)
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("Failed to parse YAML config: %v", err)
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	// Set default values
	if cfg.ReconnectDelay == "" {
		cfg.ReconnectDelay = "5s"
	}
	
	// Validate configuration
	if err := validateConfig(&cfg); err != nil {
		log.Printf("Config validation failed: %v", err)
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	log.Printf("Successfully loaded config with %d tunnels", len(cfg.Tunnels))
	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.Tunnels == nil || len(cfg.Tunnels) == 0 {
		return fmt.Errorf("no tunnels configured")
	}
	
	// Validate reconnect delay format
	if cfg.ReconnectDelay != "" {
		if _, err := time.ParseDuration(cfg.ReconnectDelay); err != nil {
			return fmt.Errorf("invalid reconnect_delay format '%s': %w", cfg.ReconnectDelay, err)
		}
	}
	
	// Track used local ports to prevent conflicts
	usedPorts := make(map[int]string)
	
	// Validate each tunnel configuration
	for name, tunnel := range cfg.Tunnels {
		if err := validateTunnel(name, tunnel, usedPorts); err != nil {
			return fmt.Errorf("tunnel '%s': %w", name, err)
		}
	}
	
	return nil
}

func validateTunnel(name string, tunnel TunnelConfig, usedPorts map[int]string) error {
	// Validate tunnel name
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("tunnel name cannot be empty")
	}
	
	// Validate host
	if strings.TrimSpace(tunnel.Host) == "" {
		return fmt.Errorf("host cannot be empty")
	}
	
	// Validate user
	if strings.TrimSpace(tunnel.User) == "" {
		return fmt.Errorf("user cannot be empty")
	}
	
	// Validate local port
	if tunnel.LocalPort <= 0 || tunnel.LocalPort > 65535 {
		return fmt.Errorf("local_port must be between 1 and 65535, got %d", tunnel.LocalPort)
	}
	
	// Check for port conflicts
	if existingTunnel, exists := usedPorts[tunnel.LocalPort]; exists {
		return fmt.Errorf("local_port %d already used by tunnel '%s'", tunnel.LocalPort, existingTunnel)
	}
	usedPorts[tunnel.LocalPort] = name
	
	// Validate remote host
	if strings.TrimSpace(tunnel.RemoteHost) == "" {
		return fmt.Errorf("remote_host cannot be empty")
	}
	
	// Validate remote port
	if tunnel.RemotePort <= 0 || tunnel.RemotePort > 65535 {
		return fmt.Errorf("remote_port must be between 1 and 65535, got %d", tunnel.RemotePort)
	}
	
	// Validate key file if specified
	if tunnel.KeyFile != "" {
		cleanKeyPath := filepath.Clean(tunnel.KeyFile)
		if strings.Contains(cleanKeyPath, "..") {
			return fmt.Errorf("invalid key file path: directory traversal detected")
		}
		
		if _, err := os.Stat(cleanKeyPath); os.IsNotExist(err) {
			return fmt.Errorf("key file does not exist: %s", cleanKeyPath)
		}
		
		// Check key file permissions (should not be world-readable)
		if info, err := os.Stat(cleanKeyPath); err == nil {
			if info.Mode().Perm()&0044 != 0 {
				log.Printf("Warning: key file %s has overly permissive permissions %o", cleanKeyPath, info.Mode().Perm())
			}
		}
	}
	
	return nil
}