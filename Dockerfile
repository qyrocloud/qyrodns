# -------- Build Stage --------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN go build -o qyrodns ./cmd/qyrodns

# -------- Runtime Stage --------
FROM alpine:3.20

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/qyrodns .

# Ensure non-root execution
RUN adduser -D -g '' qyro && chown qyro:qyro /app/qyrodns
USER qyro

# Expose required ports
EXPOSE 5300/udp 5301

# Run app
ENTRYPOINT ["./qyrodns"]
