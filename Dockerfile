# Use the official Golang image as a base
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o leetcode-to-anki-go

# Use a smaller image for the final stage
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

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