.PHONY: build run test swag docker-up docker-down

build:
	go build -o server ./cmd/api

run:
	go run ./cmd/api/main.go

test:
	go test ./... -v

swag:
	swag init -g cmd/api/main.go --output docs

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down
