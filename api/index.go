package handler

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/hadith-api/docs"
	"github.com/hadith-api/handlers"
	"github.com/hadith-api/repository"
	"github.com/hadith-api/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handler is the serverless function entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set up the repository with the data directory
	repo := repository.NewFileRepository("./data")

	// Initialize the handlers with the repository
	hadithHandler := handlers.NewHadithHandler(repo)

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Configure CORS
	router.Use(cors.Default())

	// Root redirects to Swagger UI
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Set up API routes
	apiV1 := router.Group("/api/v1")
	{
		routes.SetupHadithRoutes(apiV1, hadithHandler)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve the request
	router.ServeHTTP(w, r)
}
