# Docker Deployment Guide

## Quick Start (Zero Configuration)

The easiest deployment possible:

```bash
docker compose up -d
```

That's it! Access `https://localhost:5173` and login with:

- **Username:** admin
- **Password:** admin
- You'll be **forced to change the password** on first login

### What happens automatically:

- ✅ Pulls pre-built image from registry
- ✅ Creates secure JWT secret automatically
- ✅ Generates self-signed TLS certificates
- ✅ Creates persistent volumes for data/certs
- ✅ Runs health checks

---

## Optional Customization

### Change Initial Admin Username

If you prefer a different username instead of "admin":

```bash
# Create .env file
echo "EVOLORD_USER=myusername" > .env

docker compose up -d

# Now login with: myusername / admin
```

### Change Port

---

## Development Setup (Build from Source)

For local development:

---

## Development Setup (Build from Source)

For local development:

```bash
# Clone repository
git clone https://github.com/EvolvWTF/evolord.git
cd evolord

# Configure (leave DOCKER_IMAGE unset/commented)
cp .env.example .env
nano .env

# Build and start
docker compose build
docker compose up -d
```

---

## Publishing Images (For Maintainers)

### Automated via GitHub Actions

The repository includes a GitHub Actions workflow that automatically builds and publishes on tag creation:# Stop
docker stop evolord-server
docker rm evolord-server

````

## Build Multi-Arch Server Images

Use Docker Buildx to produce images that run on both amd64 and arm64 hosts (macOS Apple Silicon, most Linux servers):

```bash
docker buildx create --use --name evolord-builder || docker buildx use evolord-builder
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t evolord/server:multiarch \
  --load \
  .
```

## Configuration

### Environment Variables

| Variable            | Description             | Default                            |
| ------------------- | ----------------------- | ---------------------------------- |
| `EVOLORD_USER`     | Admin username          | `admin`                            |
| `EVOLORD_PASS`     | Admin password          | `admin`                            |
| `JWT_SECRET`        | JWT signing secret      | `change-this-secret-in-production` |
| `PORT`              | Server port             | `5173`                             |
| `HOST`              | Bind address            | `0.0.0.0`                          |
| `EVOLORD_TLS_CERT` | Path to TLS certificate | `./certs/server.crt`               |
| `EVOLORD_TLS_KEY`  | Path to TLS private key | `./certs/server.key`               |
| `EVOLORD_TLS_CA`   | Path to CA certificate  | (optional)                         |

### Using Custom Certificates

Mount your certificates as a volume:

```yaml
volumes:
  - ./my-certs:/app/certs:ro
```

Or set environment variables to point to your certificate paths.

### Using config.json

You can mount a `config.json` file instead of using environment variables:

```yaml
volumes:
  - ./config.json:/app/config.json:ro
```

**Note**: Environment variables take precedence over `config.json`.

## Persistence

The Docker setup persists:

- **Database**: Stored in `/app/data/evolord.db` within the `evolord-data` volume
- **Certificates**: Auto-generated on first run in `/app/certs` (ephemeral by default)

For persistent certificates across container restarts, mount the certs directory:

```yaml
volumes:
  - evolord-certs:/app/certs
```

### Backup Database

```bash
# Copy database from container
docker cp evolord-server:/app/data/evolord.db ./evolord-backup.db

# Or if using Docker Compose
docker-compose exec evolord-server cp /app/data/evolord.db /tmp/backup.db
docker cp evolord-server:/tmp/backup.db ./evolord-backup.db
```

### Restore Database

```bash
# Copy backup into container
docker cp ./evolord-backup.db evolord-server:/app/data/evolord.db
docker-compose restart evolord-server
```

## Client Configuration

Point your Evolord clients to the Docker container:

```bash
export EVOLORD_SERVER=wss://your-docker-host:5173
export EVOLORD_TLS_INSECURE_SKIP_VERIFY=true  # Development only!

cd Evolord-Client && go run ./cmd/agent
```

For production, copy the server certificate from the container:

```bash
docker cp evolord-server:/app/certs/server.crt ./server.crt
export EVOLORD_TLS_CA=./server.crt
```

## Production Deployment

### Security Checklist

- [ ] Change default username and password
- [ ] Set a strong JWT secret (minimum 32 characters)
- [ ] Use custom TLS certificates (Let's Encrypt recommended)
- [ ] Enable firewall rules
- [ ] Use reverse proxy (nginx/traefik) for additional security
- [ ] Set up regular database backups
- [ ] Review audit logs regularly

### Using with Reverse Proxy (nginx)

```nginx
upstream evolord {
    server localhost:5173;
}

server {
    listen 443 ssl http2;
    server_name evolord.example.com;

    ssl_certificate /etc/ssl/certs/your-cert.pem;
    ssl_certificate_key /etc/ssl/private/your-key.pem;

    location / {
        proxy_pass https://evolord;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Troubleshooting

### Container won't start

```bash
# Check logs
docker-compose logs evolord-server

# Common issues:
# - Port 5173 already in use
# - Invalid environment variables
# - Missing permissions for volume mounts
```

### Certificate errors

The server auto-generates self-signed certificates on first run. For production:

1. Mount custom certificates
2. Or use Let's Encrypt with a reverse proxy
3. Update clients to trust your certificate

### Database locked errors

Ensure only one container is accessing the database volume:

```bash
docker-compose down
docker-compose up -d
```

## Updates

```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## Monitoring

### Health Check

```bash
curl -k https://localhost:5173/health
```

### Resource Usage

```bash
docker stats evolord-server
```

### Audit Logs

```bash
# Access SQLite database
docker exec -it evolord-server bun run -e "
  const db = require('bun:sqlite').default(new Database('evolord.db'));
  console.log(db.query('SELECT * FROM audit_logs ORDER BY timestamp DESC LIMIT 10').all());
"
```

## Support

For issues and questions:

- Check logs: `docker-compose logs -f`
- GitHub Issues: [Your repository]
- Documentation: [README.md](README.md)
