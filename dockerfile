# -------------------------------
# 1. Build stage
# -------------------------------
FROM golang:1.25-alpine AS builder

# Metadata build stg labels
LABEL stage="builder"

# Enable Go modules as well as caching
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy our src code
COPY . .

# Build binary for linux/amd64, fully static
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -o /bin/api ./cmd/api

# -------------------------------
# 2. Final stage (minimal runtime)
# -------------------------------
FROM scratch

# Metadata labels (OCI compliant)
ARG BUILD_DATE
LABEL org.opencontainers.image.title="One Acre Fund API" \
      org.opencontainers.image.description="A production-ready Go API server built for performance and security." \
      org.opencontainers.image.authors="B.M <b.m@oneacrefund.com>" \
      org.opencontainers.image.vendor="prj-OneAcreFund" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.source="https://github.com/Blue-Davinci/one_acre_fund" \
      org.opencontainers.image.version="1.0.0" \
      org.opencontainers.image.created=$BUILD_DATE

# Set working directory
WORKDIR /app

# Copy only compiled binary
COPY --from=builder /bin/api /app/api

# Expose port (non-privileged)
EXPOSE 8080

# Run as non-root user for more security
USER 65532:65532

# No need for healthcheck, scratch file. We will use live&readiness probes

# Entrypoint
ENTRYPOINT ["/app/api"]
