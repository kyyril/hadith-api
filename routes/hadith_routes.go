package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hadith-api/handlers"
)

// SetupHadithRoutes configures all hadith-related routes
func SetupHadithRoutes(router *gin.RouterGroup, handler *handlers.HadithHandler) {
	// Get all available narrators
	router.GET("/narrators", handler.GetNarrators)
	// Get all Hadiths with pagination(limit 10 per page) and search optional
	router.GET("/hadis", handler.GetAllHadiths)
	// Get hadiths by narrator
	router.GET("/hadis/:slug", handler.GetHadithsByNarrator)
	// Get hadith by narrator and number
	router.GET("/hadis/:slug/:number", handler.GetHadithByNumber)
}
