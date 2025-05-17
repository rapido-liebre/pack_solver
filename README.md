# Pack Solver

This is a Go backend service that allows configuration of available pack sizes and calculates the optimal set of packs to fulfill an order quantity.

## Features
- Configure allowed pack sizes (e.g., 250, 500, 1000)
- Calculate pack distribution for a requested quantity
- REST API built with Gin
- Redis for storing configuration

## Requirements
- Go 1.22+
- Docker (optional for local setup)

## Run locally

### 1. Run with Docker Compose
```bash
docker-compose up --build
```

### 2. Run manually
```bash
export REDIS_ADDR=localhost:6379
redis-server &
go run ./cmd/api
```

## API

### GET /config/packs
Returns current pack sizes.

**Response:**
```json
{
  "pack_sizes": [250, 500, 1000]
}
```

---

### POST /config/packs
Sets allowed pack sizes.

**Request:**
```json
{
  "pack_sizes": [250, 500, 1000, 5000]
}
```

**Response:**
```json
{
  "success": true,
  "pack_sizes": [250, 500, 1000, 5000]
}
```

---

### POST /order
Returns pack combination to match quantity.

**Request:**
```json
{
  "quantity": 12001
}
```

**Response:**
```json
{
  "packs": [
    { "size": 5000, "count": 2 },
    { "size": 2000, "count": 1 },
    { "size": 250, "count": 1 }
  ],
  "total_items": 12000
}
```

## Project Structure
- `cmd/api` - entrypoint
- `internal/http` - route handlers
- `internal/packsolver` - pack solving logic
- `internal/config` - Redis integration
- `internal/packsolver/solver_test.go` - unit tests

## Tests
```bash
go test ./...
```

---

## License
MIT
