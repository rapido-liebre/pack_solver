// Package http internal/http/routes.go
package http

import (
	"github.com/rapido-liebre/pack_solver/internal/packsolver"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/rapido-liebre/pack_solver/internal/config"
)

type PackConfigRequest struct {
	PackSizes []int `json:"pack_sizes" binding:"required"`
}

type PackConfigResponse struct {
	Success   bool  `json:"success"`
	PackSizes []int `json:"pack_sizes"`
}

type OrderRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

type OrderResponse struct {
	Packs      []packsolver.Pack `json:"packs"`
	TotalItems int               `json:"total_items"`
}

// RegisterRoutes registers HTTP routes for managing pack size configuration:
// - GET /config/packs: returns the current pack size configuration
// - POST /config/packs: updates the pack size configuration after validation
// - POST /order: returns the optimal pack distribution for the requested quantity
func RegisterRoutes(r *gin.Engine) {
	r.GET("/config/packs", getPackSizes)
	r.POST("/config/packs", setPackSizes)
	r.POST("/order", createOrder)
}

// getPackSizes returns the current pack size configuration fetched from Redis.
// It responds with a JSON object containing the list of pack sizes.
func getPackSizes(c *gin.Context) {
	sizes, err := config.GetPackSizes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pack sizes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pack_sizes": sizes})
}

// setPackSizes sets a new pack size configuration after validation.
// It ensures all pack sizes are positive integers, removes duplicates,
// and sorts the list for consistency and solver optimization.
func setPackSizes(c *gin.Context) {
	var req PackConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.PackSizes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing pack_sizes array"})
		return
	}

	// Validation: all pack sizes must be > 0
	// This loop checks that every provided pack size is a positive integer
	for _, s := range req.PackSizes {
		if s <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "pack sizes must be > 0"})
			return
		}
	}

	// Remove duplicates and sort ascending
	sizeMap := map[int]struct{}{}
	for _, s := range req.PackSizes {
		sizeMap[s] = struct{}{}
	}
	clean := make([]int, 0, len(sizeMap))
	for s := range sizeMap {
		clean = append(clean, s)
	}

	// Sort the pack sizes in ascending order for consistency and optimization in the solver algorithm
	sort.Ints(clean)

	if err := config.SetPackSizes(clean); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not store new config"})
		return
	}

	c.JSON(http.StatusOK, PackConfigResponse{Success: true, PackSizes: clean})
}

// createOrder processes a request for an order and returns the optimal pack distribution.
func createOrder(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing quantity"})
		return
	}

	sizes, err := config.GetPackSizes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch pack sizes"})
		return
	}

	packs, total := packsolver.SolvePackDistribution(req.Quantity, sizes)
	c.JSON(http.StatusOK, OrderResponse{
		Packs:      packs,
		TotalItems: total,
	})
}
