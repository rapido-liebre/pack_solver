# 📦 Pack Solver

A Go-based backend service for calculating optimal pack combinations for a given order quantity.  
It also supports runtime configuration of available pack sizes and provides a built-in UI and Swagger documentation.

---

## 🚀 Features

- Calculate pack distributions for any quantity
- Runtime configuration of pack sizes (no code change)
- Simple HTML UI and Swagger for testing
- Redis-based persistent storage
- Dockerized with `docker-compose`
- Fully testable (unit + integration)

---

## 📦 Endpoints

### `POST /order`
Calculate optimal pack sizes for a given quantity.

Request:
```json
{ "quantity": 2300 }
```

Response:
```json
{
  "packs": [
    { "size": 1000, "count": 2 },
    { "size": 250, "count": 1 },
    { "size": 100, "count": 1 }
  ],
  "total_items": 2350
}
```

---

### `GET /config/packs`
Returns the current list of configured pack sizes.

### `POST /config/packs`
Updates the pack size configuration.

Request:
```json
{ "pack_sizes": [100, 250, 500, 1000] }
```

---

### `GET /health`
Simple health check endpoint.

---

### `GET /swagger/index.html`
Swagger UI for testing the API interactively.

---

## 🔧 Local development

### Requirements:
- Go 1.24.3+
- Redis

### Run locally:

```bash
make run
```

Swagger: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

Frontend UI: [http://localhost:8080/](http://localhost:8080/)

---

## 🧪 Run tests

```bash
make test
```

---

## 🐳 Docker

### Run using Docker Compose:

```bash
make docker-up
```

Stop containers:

```bash
make docker-down
```

---

## 📚 Generate Swagger Docs

Make sure you have [swag](https://github.com/swaggo/swag) installed:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Then run:

```bash
make swag
```

> 💡 If you encounter a build error with `docs.go` mentioning `LeftDelim` or `RightDelim`, manually remove those fields from the `SwaggerInfo` struct in `docs/docs.go`.

---

| Command             | Description                                       |
|---------------------|---------------------------------------------------|
| `make build`        | Builds the Go binary to `./server`                |
| `make run`          | Runs the application locally (`go run`)           |
| `make test`         | Runs all unit and integration tests               |
| `make swag`         | Generates Swagger documentation (`/docs`)         |
| `make docker-up`    | Builds and starts the app with Redis via Docker   |
| `make docker-down`  | Stops and removes Docker containers               |