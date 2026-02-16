# Go Links - Internal URL Shortener

A lightweight, production-ready URL shortener service for internal use. Maps short slugs (e.g., `go/wiki`) to full URLs with instant redirects.

## Features

- **Fast redirects**: GET `/slug` → 302 redirect to destination URL
- **Web UI**: Beautiful listing of all links at `/`
- **REST API**: POST `/admin/add` to create new links
- **SQLite storage**: Persistent, zero-config database
- **Basic Auth**: Optional HTTP Basic Auth for admin endpoints
- **Logging**: Request logging for all operations
- **Docker-ready**: Multi-stage build, non-root user, configurable paths

## Build Instructions

### Local Build

```bash
# Install dependencies
go mod download

# Build binary
go build -o golinks main.go

# Run
./golinks
```

### Docker Build

```bash
# Build image
docker build -t golinks:latest .

# Run container
docker run -p 8080:8080 \
  -v $(pwd)/data:/data \
  -e ADMIN_USER=admin \
  -e ADMIN_PASS=secretpass \
  golinks:latest
```

### Docker Compose

```bash
# Start service
docker-compose up -d

# View logs
docker-compose logs -f

# Stop service
docker-compose down
```

## Configuration

All configuration is via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_PATH` | `./data/links.db` | Path to SQLite database file |
| `LISTEN_ADDR` | `0.0.0.0:8080` | Server listen address and port |
| `ADMIN_USER` | _(optional)_ | Username for admin endpoints |
| `ADMIN_PASS` | _(optional)_ | Password for admin endpoints |

**Note**: If `ADMIN_USER` and `ADMIN_PASS` are not set, admin endpoints will be accessible without authentication (not recommended for production).

## API Usage

### List All Links

```bash
# Opens in browser - shows beautiful web UI
curl http://localhost:8080/
```

### Follow a Link

```bash
# Redirects to the target URL
curl -L http://localhost:8080/wiki
```

### Add a New Link

```bash
# With authentication
curl -X POST http://localhost:8080/admin/add \
  -u admin:secretpass \
  -H "Content-Type: application/json" \
  -d '{
    "slug": "wiki",
    "url": "https://wiki.company.com"
  }'

# Response
{
  "status": "created",
  "slug": "wiki",
  "url": "https://wiki.company.com"
}
```

### Remove a Link

```bash
# Remove a slug (with authentication)
curl -X POST http://localhost:8080/admin/remove \
  -u admin:secretpass \
  -H "Content-Type: application/json" \
  -d '{"slug": "wiki"}'

# Response
{
  "status": "removed",
  "slug": "wiki"
}
```
```

### Example Links

```bash
# Add common shortcuts
curl -X POST http://localhost:8080/admin/add -u admin:secretpass \
  -H "Content-Type: application/json" \
  -d '{"slug": "jira", "url": "https://jira.company.com"}'

curl -X POST http://localhost:8080/admin/add -u admin:secretpass \
  -H "Content-Type: application/json" \
  -d '{"slug": "github", "url": "https://github.com/yourorg"}'

curl -X POST http://localhost:8080/admin/add -u admin:secretpass \
  -H "Content-Type: application/json" \
  -d '{"slug": "docs", "url": "https://docs.company.com"}'
```

## Database Schema

The SQLite database has a single table:

```sql
CREATE TABLE IF NOT EXISTS links (
    slug TEXT PRIMARY KEY,
    url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## URL Validation

- Only `http://` and `https://` URLs are accepted
- URLs must be valid and parseable
- Slugs must be unique and non-empty
- Reserved slug: `admin` (cannot be used)

## Production Deployment

### Security Checklist

- ✅ Set strong `ADMIN_USER` and `ADMIN_PASS`
- ✅ Use HTTPS reverse proxy (nginx, Traefik, Caddy)
- ✅ Restrict network access to internal network only
- ✅ Regular database backups of `./data/links.db`
- ✅ Monitor logs for suspicious activity

### Reverse Proxy Example (nginx)

```nginx
server {
    listen 80;
    server_name go.company.internal;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Health Check

```bash
# Check if service is running
curl -f http://localhost:8080/ || echo "Service down"
```

## Troubleshooting

### Database Permission Errors

```bash
# Ensure data directory is writable
mkdir -p ./data
chmod 755 ./data
```

### Port Already in Use

```bash
# Change listen address
export LISTEN_ADDR=0.0.0.0:9090
./golinks
```

### View Logs

```bash
# Docker Compose
docker-compose logs -f golinks

# Standalone binary logs to stdout
./golinks 2>&1 | tee golinks.log
```

## Development

### Project Structure

```
golinks/
├── main.go              # Application code (~270 lines)
├── go.mod               # Go module definition
├── Dockerfile           # Multi-stage Docker build
├── docker-compose.yaml  # Docker Compose configuration
└── README.md            # This file
```

### Code Highlights

- **No external frameworks**: Pure `net/http` and `database/sql`
- **Clean error handling**: Proper logging and HTTP status codes
- **Modern Go practices**: Go 1.22+ idioms
- **Production-ready**: Graceful startup, validation, logging

## License

Internal use only. Not for public distribution.

## Support

For issues or feature requests, contact your infrastructure team.
