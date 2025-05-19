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

## 📦 API Endpoints

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

## 🌍 Live Demo

You can try the app live:

- 🔧 Backend UI: [https://packsolver.up.railway.app/swagger/index.html](https://packsolver.up.railway.app/swagger/index.html)
- 🖥️ Web UI: [https://packsolver.up.railway.app/](https://packsolver.up.railway.app/)


### Algorithms Used

The backend offers **three algorithms** for solving the pack distribution problem:

1. **Greedy (SolveGreedy)** – chooses the largest possible packs first and fills the remainder. Fast but not always optimal.
2. **Dynamic Programming (SolvePackDistribution)** – computes minimal excess above required amount. Optimal but slower for very large input.
3. **Smart Strategy (SolveSmart)** – runs both Greedy and DP and picks the better result based on the lowest total amount.

The `/order` endpoint uses the Smart strategy by default.

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

## ⚙️ Environment configuration

To run the project locally, create a `.env` file based on the `.env.sample` provided.

Example `.env`:

```
REDIS_ADDR=localhost:6379
PACK_SOLVER_API=http://localhost:8080
```

The project uses `github.com/joho/godotenv` to load variables automatically.

---

## 🧪 Run tests

Unit and integration tests:

```bash
make test
```
E2E test (requires local server running on port 8080 and Redis):
```bash
go test tests/e2e_test.go
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

---

## 🚀 Deployment

### 🧠 Backend (Railway)

1. Log in at [https://railway.app](https://railway.app)
2. Create a new project and link your GitHub repository
3. Railway will auto-detect the Dockerfile and build the service
4. Add environment variable:
   ```
   REDIS_ADDR = ${{Redis.REDIS_URL}}
   ```
5. Redis service is provisioned automatically (defined in `railway.toml`)
6. Access your backend via Railway's generated domain

### 🌐 Frontend (Optional: Netlify or GitHub Pages)

1. Copy the contents of `/ui` into a separate repository (e.g. `pack-solver-ui`)
2. Deploy it as a static site using Netlify or GitHub Pages
3. If needed, configure CORS in the backend to allow requests from the frontend origin