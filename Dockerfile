# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate GraphQL code
RUN go run github.com/99designs/gqlgen generate

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o bin/server ./cmd/server/main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates, netcat (for health checks) and apply security updates
RUN apk update && \
    apk --no-cache add ca-certificates netcat-openbsd tzdata && \
    rm -rf /var/cache/apk/*

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/bin/server .

# Copy migrations directory
COPY ./migrations ./migrations

# Copy entrypoint script
COPY docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

# Set entrypoint
ENTRYPOINT ["/docker-entrypoint.sh"]

# Run the application
CMD ["./server"]

