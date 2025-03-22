FROM golang:1.23-alpine

WORKDIR /app

# Copy go.mod first
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main ./cmd/main.go

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./main"] 