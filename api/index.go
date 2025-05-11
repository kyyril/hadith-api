// api/index.go
package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hadith-api/config"
	_ "github.com/hadith-api/docs" // This will import the docs package and run its init() function
	"github.com/hadith-api/handlers"
	"github.com/hadith-api/repository"
	"github.com/hadith-api/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handler is the serverless function entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set GO_ENV to production for Vercel
	os.Setenv("GO_ENV", "production")

	// Set BASE_URL environment variable for Swagger if not already set
	if os.Getenv("BASE_URL") == "" {
		os.Setenv("BASE_URL", "https://hadith-api-go.vercel.app")
	}

	// Get configuration after setting environment variables
	cfg := config.GetConfig()

	// Log current configuration
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("BaseURL: %s", cfg.BaseURL)

	// Continue with the original logic
	cwd, _ := os.Getwd()

	// Try multiple possible data directory locations
	possiblePaths := []string{
		filepath.Join(cwd, "data"),        // /var/task/data
		"/var/task/data",                  // Direct reference
		filepath.Join(cwd, "..", "data"),  // One level up
		filepath.Join(cwd, "api", "data"), // In api subdirectory
		"./data",                          // Relative path
		"../data",                         // Relative path one level up
	}

	var dataDir string
	var foundValidPath bool

	// Log all paths we're checking
	log.Printf("Current working directory: %s", cwd)
	log.Printf("Checking for data directory in multiple locations...")

	for _, path := range possiblePaths {
		log.Printf("Checking path: %s", path)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			log.Printf("Found data directory at: %s", path)
			dataDir = path
			foundValidPath = true
			break
		}
	}

	if !foundValidPath {
		// List contents of current directory to help debug
		files, _ := os.ReadDir(cwd)
		fileList := "Directory contents of " + cwd + ": "
		for _, file := range files {
			fileList += file.Name() + ", "
		}
		errorMsg := fmt.Sprintf("Data directory not found in any expected location. %s", fileList)
		log.Printf(errorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMsg))
		return
	}

	// Try to list JSON files in the data directory
	files, err := os.ReadDir(dataDir)
	if err != nil {
		errorMsg := fmt.Sprintf("Found data directory but could not read files: %v", err)
		log.Printf(errorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMsg))
		return
	}

	jsonFiles := []string{}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			jsonFiles = append(jsonFiles, file.Name())
		}
	}

	if len(jsonFiles) == 0 {
		errorMsg := "No JSON files found in data directory: " + dataDir
		log.Printf(errorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMsg))
		return
	}

	// Set up the repository with the determined data directory
	repo := repository.NewFileRepository(dataDir)

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

	// Health check endpoint with more detailed info
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":     "ok",
			"workingDir": cwd,
			"dataDir":    dataDir,
			"jsonFiles":  jsonFiles,
			"baseURL":    cfg.BaseURL,
			"env":        cfg.Environment,
		})
	})

	// Set up API routes
	apiV1 := router.Group("/api/v1")
	{
		routes.SetupHadithRoutes(apiV1, hadithHandler)
	}

	// Swagger documentation with URL customization
	url := ginSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", cfg.BaseURL))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Serve the request
	router.ServeHTTP(w, r)
}
