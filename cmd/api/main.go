// cmd/api/main.go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rapido-liebre/pack_solver/internal/config"
	"github.com/rapido-liebre/pack_solver/internal/http"
)

func main() {
	if err := config.InitRedis(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	r := gin.Default()
	http.RegisterRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
