package main

import (
	"log"
	"os"

	"flash-sale-coupon-system/internal/config"
	"flash-sale-coupon-system/internal/database"
	"flash-sale-coupon-system/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize handlers
	couponHandler := handlers.NewCouponHandler(db)

	// Setup router
	router := gin.Default()

	// API routes
	api := router.Group("/api")
	{
		coupons := api.Group("/coupons")
		{
			coupons.POST("", couponHandler.CreateCoupon)
			coupons.POST("/claim", couponHandler.ClaimCoupon)
			coupons.GET("/:name", couponHandler.GetCouponDetails)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
