# Build stage
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Install build dependencies and security updates
RUN apt-get update && \
    apt-get install -y git gcc libc6-dev && \
    apt-get upgrade -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy go mod files
COPY go.mod ./
COPY go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Tidy modules and build the application
RUN go mod tidy && CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:3.20

# Install runtime dependencies and security updates
RUN apk update && \
    apk --no-cache add ca-certificates sqlite wget && \
    apk upgrade --no-cache && \
    rm -rf /var/cache/apk/*

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Create app directory
RUN mkdir -p /var/app && chown -R appuser:appgroup /var/app

# Copy binary from builder stage
COPY --from=builder /app/main /usr/local/bin/newsletter

# Copy static files (SvelteKit build output)
COPY --from=builder /app/static /var/app/static

# Copy migrations
COPY --from=builder /app/migrations /var/app/migrations

# Set working directory
WORKDIR /var/app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
CMD ["newsletter"]
