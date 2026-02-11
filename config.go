package main

import (
	"os"
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
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.ReconnectDelay == "" {
		cfg.ReconnectDelay = "5s"
	}
	return &cfg, nil
}
