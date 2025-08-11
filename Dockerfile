# Build stage
FROM golang:1.23.5-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates (needed for fetching dependencies and HTTPS)
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Production stage
FROM alpine:latest

# Install ca-certificates, timezone data, and netcat for health checks
RUN apk --no-cache add ca-certificates netcat-openbsd tzdata

# Set timezone (optional, adjust as needed)
ENV TZ=Asia/Jakarta

# Create app user
RUN addgroup -g 1001 -S appuser && \
    adduser -S -D -H -u 1001 -s /sbin/nologin appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Create uploads directory with proper structure and permissions
RUN mkdir -p uploads/news uploads/members && \
    chown -R appuser:appuser /app && \
    chmod -R 755 /app/uploads

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD nc -z localhost 8080 || exit 1

# Run the application
CMD ["./main"]
