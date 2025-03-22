# Blueis - A Redis Clone

Blueis is a lightweight Redis clone written in Go. It implements the Redis Serialization Protocol (RESP) and provides basic Redis functionality.

## Features

- RESP protocol support
- Basic Redis commands:
  - SET
  - GET
  - DEL
  - PING
- Thread-safe operations
- Docker support

## Getting Started

### Prerequisites

- Go 1.23 or later
- Docker (optional)

### Running Locally

1. Clone the repository:
```bash
git clone https://github.com/yourusername/blueis.git
cd blueis
```

2. Run the server:
```bash
go run cmd/main.go
```

The server will start on port 8080.

### Running with Docker

1. Build and run using docker-compose:
```bash
docker-compose up --build
```

2. Or build and run using Docker directly:
```bash
docker build -t blueis .
docker run -p 8080:8080 blueis
```

## Testing

You can test the server using the Redis CLI:

```bash
redis-cli -p 8080
```

Example commands:
```redis
SET mykey "Hello World"
GET mykey
PING
DEL mykey
```

## Project Structure

```
.
├── cmd/
│   └── main.go         # Application entry point
├── internal/
│   ├── redis/         # Redis implementation
│   │   └── redis.go
│   ├── resp/          # RESP protocol implementation
│   │   └── resp.go
│   └── server/        # TCP server implementation
│       └── server.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 