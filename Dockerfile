# ---------- BUILD STAGE ----------
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# ---------- RUNTIME STAGE ----------
FROM scratch

WORKDIR /root/

# Copy binary and assets
COPY --from=builder /app/server .
COPY --from=builder /app/ui ./ui
COPY --from=builder /app/docs ./docs

# Run binary
ENTRYPOINT ["./server"]
