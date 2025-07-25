# ─── Builder Stage ────────────────────────────────────────────────────────────
FROM golang:1.24.2-alpine AS builder

#? Install git (for private modules) and set working directory
#? RUN apk add --no-cache git
WORKDIR /app

# Cache dependencies for faster rebuilds
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build the 'web' binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w" \
    -o cards-site ./cmd/web

# ─── Final Stage ──────────────────────────────────────────────────────────────
FROM alpine:latest

# Include certificates for HTTPS clients
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy only the compiled binary from builder
COPY --from=builder /app/cards-site ./

# For the image to include config.yml without mounting:
COPY --from=builder /app/cmd/web/config.yaml cmd/web/config.yaml

# Executable
RUN chmod +x cards-site

# Expose the port defined in your config.yml
EXPOSE 8080

# Default entrypoint and arguments
ENTRYPOINT ["./cards-site"]
CMD ["--config", "cmd/web/config.yaml"]