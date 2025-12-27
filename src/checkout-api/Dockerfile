# --- Stage 1: Builder ---
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary statically linked
# CGO_ENABLED=0 is crucial for running in distroless/scratch
RUN CGO_ENABLED=0 GOOS=linux go build -o checkout-api .

# --- Stage 2: Runner ---
# Use Google's distroless static image (no shell, no package manager)
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy binary from builder stage
COPY --from=builder /app/checkout-api /checkout-api

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/checkout-api"]