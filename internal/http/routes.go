// Package http internal/http/routes.go
package http

import (
	"github.com/rapido-liebre/pack_solver/internal/packsolver"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/rapido-liebre/pack_solver/internal/config"

	_ "github.com/rapido-liebre/pack_solver/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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
	Packs      []packsolver.PackResult `json:"packs"`
	TotalItems int                     `json:"total_items"`
}

// SetupRouter initializes the Gin engine with all registered routes.
// Used by main() and tests to start the API server.
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	RegisterRoutes(r)
	return r
}

// RegisterRoutes registers HTTP routes for managing pack size configuration and serving frontend UI.
// It also exposes Swagger documentation and a health check endpoint.
// - GET /: serve index.html as default
// - GET /config/packs: returns the current pack size configuration
// - POST /config/packs: updates the pack size configuration after validation
// - POST /order: returns the optimal pack distribution for the requested quantity
func RegisterRoutes(r *gin.Engine) {
	// Serve UI from /ui directory
	r.Static("/static", "./ui")
	r.GET("/", func(c *gin.Context) {
		c.File("./ui/index.html")
	})

	r.GET("/config/packs", getPackSizes)
	r.POST("/config/packs", setPackSizes)
	r.POST("/order", createOrder)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// @Summary Get current pack size configuration
// @Description Returns the list of configured pack sizes fetched from Redis
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]int
// @Failure 500 {object} map[string]string
// @Router /config/packs [get]
func getPackSizes(c *gin.Context) {
	sizes, err := config.GetPackSizes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pack sizes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pack_sizes": sizes})
}

// @Summary Update pack size configuration
// @Description Set a new list of pack sizes (must be unique and > 0). It ensures all pack sizes are positive integers, removes duplicates,
// and sorts the list for consistency and solver optimization
// @Tags config
// @Accept json
// @Produce json
// @Param request body PackConfigRequest true "Pack sizes"
// @Success 200 {object} PackConfigResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /config/packs [post]
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

// @Summary Calculate pack distribution
// @Description Calculates the optimal pack combination for the requested quantity
// @Tags order
// @Accept json
// @Produce json
// @Param request body OrderRequest true "Order quantity"
// @Success 200 {object} OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /order [post]
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
