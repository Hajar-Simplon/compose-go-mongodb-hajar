# Docker Compose - Go & MongoDB

## Project Overview

This project demonstrates how to use **Docker Compose** to orchestrate a multi-container application composed of:

- A Go web application
- A MongoDB database
- An Nginx reverse proxy / load balancer

---

## Architecture

```text
                Browser
                    │
         http://localhost:5000
                    │
                    ▼
      Nginx Reverse Proxy / Load Balancer
                    │
        ┌───────────┼───────────┐
        ▼           ▼           ▼
     Go App      Go App      Go App
                    │
                    ▼
                MongoDB
                    │
                    ▼
             Docker Volume
         (Persistent Storage)
```

---

## Features

- Multi-container application managed by Docker Compose
- Automatic image build for the Go application
- MongoDB official Docker image
- Dedicated Docker Compose network
- Application exposed on port **5000**
- Persistent MongoDB storage using Docker Volumes
- Automatic container restart
- MongoDB health check
- Go application starts only after MongoDB is healthy
- Docker Compose Watch support for automatic rebuild during development
- Docker Secrets for secure management of sensitive data
- Nginx reverse proxy
- Horizontal scaling with Docker Compose
- Dynamic DNS resolution for load balancing
- Traffic distribution across multiple application replicas

---

## Project Structure

```text
.
├── Dockerfile
├── compose.yml
├── go.mod
├── go.sum
├── handlers.go
├── main.go
├── Makefile
├── nginx/
│   ├── Dockerfile
│   └── default.conf
├── secrets/
│   └── mongodb_password.txt
├── .gitignore
└── README.md
```

---

## Prerequisites

- Docker Engine 29+
- Docker Compose V2

Verify installation:

```bash
docker --version
docker compose version
```

---

## Start the Application

Build and start all services:

```bash
docker compose up --build
```

Or in detached mode:

```bash
docker compose up -d --build
```

---

## Stop the Application

```bash
docker compose down
```

To also delete the persistent database:

```bash
docker compose down -v
```

---

## Access the Application

Application:

```
http://localhost:5000
```

Logs endpoint:

```
http://localhost:5000/logs
```

---

## Docker Compose Features

### Build

The Go application image is built automatically using the local Dockerfile.

```yaml
build: .
```

---

### Networking

Docker Compose automatically creates an isolated network.

The Go application connects to MongoDB using the service name:

```
mongodb
```

instead of using an IP address.

---

### Persistent Storage

MongoDB stores its database inside:

```
/data/db
```

A named Docker volume is mounted on this directory:

```yaml
volumes:
  - mongodb-data:/data/db
```

This allows the database to persist even after containers are removed.

---

### Restart Policy

Both containers use:

```yaml
restart: unless-stopped
```

Containers automatically restart after a failure unless they were manually stopped.

---

### Health Check

MongoDB is periodically checked using:

```yaml
healthcheck:
  test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
```

The Go application only starts once MongoDB is healthy.

---

### Docker Compose Watch

Docker Compose watches the project directory.

Whenever the application source changes, Docker automatically rebuilds the image and restarts the container.

Start watch mode:

```bash
docker compose up --watch
```

---

### Docker Secrets

Docker Secrets provide a secure way to manage sensitive information such as passwords or API keys.

Instead of storing sensitive values directly in the `compose.yml` file or committing them to Git, Docker mounts them inside the container as files located in:

```
/run/secrets/
```

In this project, a MongoDB password is stored as a Docker Secret.

```yaml
secrets:
  mongodb_password:
    file: ./secrets/mongodb_password.txt
```

The secret is mounted into the MongoDB container:

```yaml
mongodb:
  secrets:
    - mongodb_password
```

The secret file is excluded from version control using `.gitignore`.

---
## Load Balancing & Scaling

### Objective

To improve the application's capacity, the Go service can be scaled horizontally by running multiple replicas.

Instead of exposing the Go application directly, an Nginx reverse proxy becomes the single public entry point and distributes incoming requests between the application replicas.

---

### Updated Architecture

```text
                Browser
                    │
         http://localhost:5000
                    │
                    ▼
        Nginx Reverse Proxy
          │       │       │
          ▼       ▼       ▼
       Go App  Go App  Go App
              │
              ▼
           MongoDB
```

---

### Scaling the Application

Start three replicas of the Go application:

```bash
docker compose up -d --scale app=3
```

Docker Compose automatically creates:

- app-1
- app-2
- app-3

---

### Nginx Reverse Proxy

The Go application is no longer exposed directly to the host.

Instead of:

```yaml
ports:
  - "5000:8080"
```

the application now uses:

```yaml
expose:
  - "8080"
```

Only the Nginx service exposes a port:

```yaml
ports:
  - "5000:80"
```

This makes Nginx the single entry point for all incoming requests.

---

### Dynamic DNS Resolution

Initially, the Nginx configuration used:

```nginx
proxy_pass http://app:8080/;
```

Although Docker Compose successfully created three application replicas, all requests were forwarded to the same container.

Investigation showed that Docker's embedded DNS correctly returned the IP addresses of all replicas, but Nginx resolved the hostname only once during startup and cached the first IP address.

The solution was to configure Docker's embedded DNS resolver.

```nginx
resolver 127.0.0.11 valid=5s ipv6=off;

location / {

    set $backend http://app:8080;

    proxy_pass $backend;

}
```

This configuration forces Nginx to resolve the backend dynamically instead of permanently caching a single IP address, allowing requests to be distributed across all application replicas.

---

### Verification

Scale the application:

```bash
docker compose up -d --scale app=3
```

Monitor the application logs:

```bash
docker compose logs -f app
```

Refresh:

```
http://localhost:5000
```

Requests should now be handled by all three replicas, confirming that load balancing is working correctly.

---

## Useful Commands

Start services

```bash
docker compose up -d
```

Rebuild services

```bash
docker compose up --build
```

View logs

```bash
docker compose logs
```

Follow logs

```bash
docker compose logs -f
```

List containers

```bash
docker compose ps
```

List volumes

```bash
docker volume ls
```

Inspect Compose configuration

```bash
docker compose config
```

---

## Verifying Persistent Storage

1. Open:

```
http://localhost:5000
```

2. Check stored documents:

```
http://localhost:5000/logs
```

3. Stop the containers:

```bash
docker compose down
```

4. Restart:

```bash
docker compose up -d
```

5. Visit:

```
http://localhost:5000/logs
```

The data should still be present because it is stored inside the Docker volume.

---

## Technologies

- Go
- MongoDB
- Docker
- Docker Compose

---

## Author

Hajar Mezouar