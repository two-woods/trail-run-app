package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// TODO: Load configuration from environment variables
	// TODO: Initialize database connection (GORM)
	// TODO: Initialize Redis connection
	// TODO: Setup JWT middleware
	// TODO: Register routes

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", registerHandler)
			auth.POST("/login", loginHandler)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(authMiddleware())
		{
			// User routes
			protected.GET("/users/me", getCurrentUserHandler)

			// Race routes
			protected.GET("/races", listRacesHandler)
			protected.GET("/races/:id", getRaceHandler)
			protected.GET("/races/:id/plan", getRacePlanHandler)
			protected.GET("/races/:id/equipment", getEquipmentHandler)
			protected.POST("/races/:id/equipment/check", checkEquipmentHandler)

			// Results routes
			protected.GET("/results", listResultsHandler)
			protected.POST("/results", createResultHandler)
			protected.GET("/results/:id", getResultHandler)
			protected.POST("/results/:id/image", generateImageHandler)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Placeholder handlers - to be implemented

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement JWT validation
		c.Next()
	}
}

func registerHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func loginHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func getCurrentUserHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func listRacesHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func getRaceHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func getRacePlanHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func getEquipmentHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func checkEquipmentHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func listResultsHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func createResultHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func getResultHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func generateImageHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
