package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hadith-api/handlers"
	"github.com/hadith-api/repository"
	"github.com/hadith-api/routes"

	_ "github.com/hadith-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Hadith API
// @version 1.0
// @description API Hadis dengan terjemahan dari 9 perawi
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.basic BasicAuth
func main() {
	// Set up the repository with the data directory
	repo := repository.NewFileRepository("./data")

	// Initialize the handlers with the repository
	hadithHandler := handlers.NewHadithHandler(repo)

	// Set up Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.Default())

	// Root redirects to Swagger UI
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Set up API routes
	apiV1 := r.Group("/api/v1")
	{
		routes.SetupHadithRoutes(apiV1, hadithHandler)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
