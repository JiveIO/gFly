# Docker Configuration

This directory contains Docker configuration files for setting up the development environment.

## Purpose

The docker directory is used for:
- Defining Docker services required by the application
- Providing configuration files for each service
- Enabling consistent development environments
- Simplifying the setup process for new developers

## Structure

- **docker-compose.yml**: Main Docker Compose configuration file
- **mailpit/**: Configuration files for the Mailpit service
  - **README.md**: Instructions for generating TLS certificates
  - **authfile**: Authentication configuration for Mailpit
  - **cert.pem**: TLS certificate for secure connections
  - **key.pem**: TLS private key for secure connections
- **redis/**: Configuration files for the Redis service
  - **redis.env**: Environment variables for Redis configuration

## Services

The Docker Compose configuration includes the following services:

### PostgreSQL Database (db)
- Image: postgres:16.4-alpine
- Port: 5432
- Credentials: user/secret
- Database: gfly

### Mail Testing (mail)
- Image: axllent/mailpit
- Ports: 8025 (web interface), 1025 (SMTP)
- Web interface: http://localhost:8025
- SMTP server: localhost:1025

### Redis (redis)
- Image: redis:7.4.0-alpine
- Port: 6379
- Password: secret (defined in redis.env)

## Usage

To start all services:

```bash
make docker.run
```

To stop all services:

```bash
make docker.stop
```

To view logs:

```bash
docker logs gfly-db    # Database logs
docker logs gfly-mail  # Mail logs
docker logs gfly-redis # Redis logs
```

## Best Practices

- Do not commit sensitive information in Docker configuration files
- Use environment variables for secrets and configuration
- Keep Docker images updated to the latest stable versions
- Use health checks to ensure services are properly initialized
- Document any special setup requirements for each service
- Use volumes for persistent data storage
