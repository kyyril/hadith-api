package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hadith-api/models"
	"github.com/hadith-api/repository"
)

// HadithHandler handles HTTP requests related to hadiths
type HadithHandler struct {
	repo *repository.FileRepository
}

// NewHadithHandler creates a new HadithHandler with the given repository
func NewHadithHandler(repo *repository.FileRepository) *HadithHandler {
	return &HadithHandler{
		repo: repo,
	}
}

// GetNarrators godoc
// @Summary      Get list of available narrators
// @Description  Returns a list of all available hadith narrators
// @Tags         narrators
// @Produce      json
// @Success      200  {object}  models.HadithResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /narrators [get]
func (h *HadithHandler) GetNarrators(c *gin.Context) {
	narrators, err := h.repo.GetAvailableNarrators()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to get narrators",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.HadithResponse{
		Status:  "success",
		Message: "Narrators retrieved successfully",
		Data: models.Narrators{
			Available: narrators,
		},
	})
}

// GetHadithsByNarrator godoc
// @Summary      Get hadiths by narrator
// @Description  Returns all hadiths from a specific narrator with optional pagination and filtering
// @Tags         hadiths
// @Produce      json
// @Param        slug   path      string  true  "Narrator slug (e.g., muslim, bukhari)"
// @Param        page   query     int     false "Page number for pagination"
// @Param        limit  query     int     false "Items per page for pagination"
// @Param        q      query     string  false "Search query to filter hadiths"
// @Success      200    {object}  models.PaginatedResponse
// @Failure      404    {object}  models.ErrorResponse
// @Failure      500    {object}  models.ErrorResponse
// @Router       /hadis/{slug} [get]
func (h *HadithHandler) GetHadithsByNarrator(c *gin.Context) {
	narrator := c.Param("slug")
	
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	query := c.Query("q")

	// Set default pagination values
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get hadiths with pagination and filtering
	hadiths, totalItems, err := h.repo.GetHadithsByNarrator(narrator, models.QueryParams{
		Page:  page,
		Limit: limit,
		Query: query,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to get hadiths",
			Error:   err.Error(),
		})
		return
	}

	// Calculate pagination values
	totalPages := (totalItems + limit - 1) / limit // Ceiling division

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Status:  "success",
		Message: "Hadiths retrieved successfully",
		Data:    hadiths,
		Pagination: models.Pagination{
			CurrentPage: page,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			PerPage:     limit,
		},
	})
}

// GetHadithByNumber godoc
// @Summary      Get hadith by narrator and number
// @Description  Returns a specific hadith from a narrator by its number
// @Tags         hadiths
// @Produce      json
// @Param        slug    path      string  true  "Narrator slug (e.g., muslim, bukhari)"
// @Param        number  path      int     true  "Hadith number"
// @Success      200     {object}  models.HadithResponse
// @Failure      404     {object}  models.ErrorResponse
// @Failure      500     {object}  models.ErrorResponse
// @Router       /hadis/{slug}/{number} [get]
func (h *HadithHandler) GetHadithByNumber(c *gin.Context) {
	narrator := c.Param("slug")
	numberStr := c.Param("number")

	// Convert number string to int
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid hadith number",
			Error:   "Hadith number must be an integer",
		})
		return
	}

	// Get the hadith
	hadith, err := h.repo.GetHadithByNumber(narrator, number)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status:  "error",
			Message: "Hadith not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.HadithResponse{
		Status:  "success",
		Message: "Hadith retrieved successfully",
		Data:    hadith,
	})
}