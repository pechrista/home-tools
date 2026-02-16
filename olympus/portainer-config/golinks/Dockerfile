# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy go mod and source files
COPY go.mod main.go ./

# Download dependencies and generate go.sum based on imports
RUN go mod tidy && go mod download && go mod verify

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w' -o golinks .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates su-exec

# Create non-root user
RUN addgroup -g 1000 golinks && \
    adduser -D -u 1000 -G golinks golinks

# Create data directory
RUN mkdir -p /data && chown golinks:golinks /data

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/golinks .

# Add entrypoint to prepare runtime (ensure /data exists and is writable)
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Switch to non-root user
USER root

# Set default environment variables
ENV DB_PATH=/data/links.db
ENV LISTEN_ADDR=0.0.0.0:8080

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./golinks"]
