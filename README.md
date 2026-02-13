# portpilot

SSH tunnel manager. Define tunnels in YAML, auto-reconnect.

## Install

```bash
git clone https://github.com/openkickstartai/portpilot.git
cd portpilot && go build -o portpilot .
```

## Usage

```bash
# Start all tunnels from config
portpilot up

# Start specific tunnel
portpilot up db-prod

# Show tunnel status
portpilot status

# Stop all
portpilot down

# Use a custom config file
portpilot -c /path/to/tunnels.yml up

# Or via environment variable
export PORTPILOT_CONFIG=/path/to/tunnels.yml
portpilot up
```


## Config (~/.portpilot.yml)

```yaml
tunnels:
  db-prod:
    host: bastion.company.com
    user: deploy
    local_port: 5433
    remote_host: db.internal
    remote_port: 5432
    key: ~/.ssh/id_ed25519
  
  redis-staging:
    host: jump.staging.com
    user: dev
    local_port: 6380
    remote_host: redis.internal
    remote_port: 6379

auto_reconnect: true
reconnect_delay: 5s
```

## Testing

```bash
go test -v ./...
```
