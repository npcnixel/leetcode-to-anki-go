FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled
ENV CGO_ENABLED=1
RUN go build -o leetcode-to-anki-go

# Use debian:stable-slim for the runtime image
FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates libc6 && \
    rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/leetcode-to-anki-go .

# Create necessary directories for mounting
RUN mkdir -p /app/input /app/output

# Set volume mounts for input and output directories
VOLUME ["/app/input", "/app/output"]

# This application doesn't expose any ports

# Set the entrypoint to run the application
ENTRYPOINT ["/app/leetcode-to-anki-go"]

# Default command (can be overridden)
CMD []

# Label with correct project name
LABEL org.opencontainers.image.title="leetcode-to-anki-go"
LABEL org.opencontainers.image.description="LeetCode to Anki card converter written in Go" 