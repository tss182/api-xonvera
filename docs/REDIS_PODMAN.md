# Redis on Podman - Quick Reference

## Installation & Setup

### 1. Initialize Podman (first time only)
```bash
podman machine init
podman machine start
```

### 2. Run Redis Container
```bash
podman run -d --name redis -p 6379:6379 redis:latest
```

### 3. Verify Redis is Running
```bash
podman ps | grep redis
podman exec redis redis-cli PING
# Output: PONG
```

## Common Commands

### Check Redis Status
```bash
podman exec redis redis-cli INFO
podman exec redis redis-cli DBSIZE
```

### Connect to Redis CLI
```bash
podman exec -it redis redis-cli
```

### Stop Redis
```bash
podman stop redis
```

### Start Redis (after stopping)
```bash
podman start redis
```

### Remove Redis Container
```bash
podman rm redis  # Must stop first if running
```

### View Logs
```bash
podman logs redis
podman logs -f redis  # Follow logs
```

## Running Your App with Redis

```bash
# Terminal 1: Ensure Redis is running
podman ps | grep redis

# Terminal 2: Run your app
make dev
# Or
go run ./cmd/main.go
```

You should see:
```
Connected to Redis {"host": "localhost", "port": "6379"}
```

## Redis Configuration

Your app uses these environment variables:
```dotenv
REDIS_HOST=localhost       # Podman forwards to localhost
REDIS_PORT=6379
REDIS_PASSWORD=            # Empty for no password
REDIS_DB=0
```

## Helpful Podman Commands

### Podman Machine Management
```bash
podman machine ls              # List machines
podman machine start           # Start machine
podman machine stop            # Stop machine
podman machine rm              # Remove machine
```

### Container Management
```bash
podman ps                      # List running containers
podman ps -a                   # List all containers
podman restart redis           # Restart Redis
podman stats redis             # View Redis resource usage
```

### Debugging
```bash
podman inspect redis           # Get detailed container info
podman logs redis -n 100       # Last 100 log lines
podman exec redis redis-cli CONFIG GET "*"  # Redis config
```

## Data Persistence

To persist Redis data across restarts:

```bash
podman run -d \
  --name redis \
  -p 6379:6379 \
  -v redis-data:/data \
  redis:latest redis-server --appendonly yes
```

## Cleanup

To completely reset:
```bash
podman stop redis
podman rm redis
podman volume rm redis-data  # If using volumes
```

## Troubleshooting

### Redis Won't Start
```bash
podman logs redis
```

### Can't Connect from App
- Ensure Redis is running: `podman ps | grep redis`
- Check port mapping: `podman port redis`
- Verify connection: `podman exec redis redis-cli PING`

### Performance Issues
```bash
# Check Redis memory usage
podman exec redis redis-cli INFO memory

# Monitor commands in real-time
podman exec -it redis redis-cli MONITOR
```
