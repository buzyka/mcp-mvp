# Stage 1: build the Go binary
FROM golang:1.24-alpine AS builder

# Install git (needed for module fetch)
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod, go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o mcp-proxy main.go

# Stage 2: create a minimal image
FROM scratch

# Copy the compiled binary
COPY --from=builder /app/mcp-proxy /usr/local/bin/mcp-proxy

# The service listens on 8089
EXPOSE 8089

# Default entrypoint
ENTRYPOINT ["/usr/local/bin/mcp-proxy"]
