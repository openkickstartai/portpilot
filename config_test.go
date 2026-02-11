package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tmp := t.TempDir()
	cfg := filepath.Join(tmp, "config.yml")
	os.WriteFile(cfg, []byte("tunnels:\n  db:\n    host: bastion.com\n    user: deploy\n    local_port: 5433\n    remote_host: db.internal\n    remote_port: 5432\nauto_reconnect: true\n"), 0644)
	c, err := LoadConfig(cfg)
	if err != nil { t.Fatal(err) }
	if len(c.Tunnels) != 1 { t.Errorf("want 1 tunnel, got %d", len(c.Tunnels)) }
	db := c.Tunnels["db"]
	if db.Host != "bastion.com" { t.Errorf("want bastion.com, got %s", db.Host) }
	if db.LocalPort != 5433 { t.Errorf("want 5433, got %d", db.LocalPort) }
	if !c.AutoReconnect { t.Error("want auto_reconnect true") }
	if c.ReconnectDelay != "5s" { t.Errorf("want default 5s, got %s", c.ReconnectDelay) }
}

func TestLoadConfigMissing(t *testing.T) {
	_, err := LoadConfig("/nonexistent")
	if err == nil { t.Error("want error") }
}
