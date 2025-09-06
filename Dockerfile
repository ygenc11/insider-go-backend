# Multi-stage build for insider-go-backend
# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install build tools (if needed)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary (static)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

# Runtime stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app

# Non-root user (distroless already defines nonroot uid 65532)
USER 65532:65532

COPY --from=builder /out/api /app/api
COPY .env /app/.env

# Create logs directory writable by nonroot (distroless lacks shell; keep simple)
# We'll rely on mounted volume from host for logs.
EXPOSE 8080

ENV PORT=8080

ENTRYPOINT ["/app/api"]
