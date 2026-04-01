package main

import (
	"os"

	"trail-run-app/internal/handler"
	"trail-run-app/internal/model"
	"trail-run-app/pkg/middleware"
	"trail-run-app/pkg/redis"
	"trail-run-app/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	utils.InitLogger()
	utils.LogInfo("Starting Trail Run App server...")

	// Initialize JWT
	middleware.InitJWT()

	// Initialize database
	if err := model.InitDB(); err != nil {
		utils.LogError("Failed to initialize database: %v", err)
		os.Exit(1)
	}

	// Auto migrate database schema
	if err := model.AutoMigrate(); err != nil {
		utils.LogError("Failed to migrate database: %v", err)
		os.Exit(1)
	}

	// Initialize Redis
	if err := redis.InitRedis(); err != nil {
		utils.LogWarn("Warning: Redis connection failed: %v (continuing without Redis)", err)
	}

	// Setup Gin router
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		utils.RespondSuccess(c, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthRequired())
		{
			// User routes
			protected.GET("/users/me", handler.GetCurrentUser)

			// Race routes
			protected.GET("/races", handler.ListRaces)
			protected.GET("/races/:id", handler.GetRace)
			protected.GET("/races/:id/plan", handler.GetRacePlan)
			protected.GET("/races/:id/equipment", handler.GetEquipment)
			protected.POST("/races/:id/equipment/check", handler.CheckEquipment)

			// Results routes
			protected.GET("/results", handler.ListResults)
			protected.POST("/results", handler.CreateResult)
			protected.GET("/results/:id", handler.GetResult)
			protected.POST("/results/:id/image", handler.GenerateImage)

			// Hotel routes (for race day planning)
			protected.GET("/hotels/search", handler.SearchHotels)
			protected.GET("/hotels/favorites", handler.ListFavoriteHotels)
			protected.POST("/hotels/favorites", handler.CreateFavoriteHotel)
			protected.DELETE("/hotels/favorites/:id", handler.DeleteFavoriteHotel)
			protected.GET("/races/:id/hotels", handler.GetAllHotelsForRace)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.LogInfo("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		utils.LogError("Failed to start server: %v", err)
		os.Exit(1)
	}
}
