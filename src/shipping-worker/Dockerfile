# --- Stage 1: Builder ---
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -o shipping-worker .

# --- Stage 2: Runner ---
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy binary from builder stage
COPY --from=builder /app/shipping-worker /shipping-worker

# No EXPOSE needed for a background worker

USER nonroot:nonroot

ENTRYPOINT ["/shipping-worker"]