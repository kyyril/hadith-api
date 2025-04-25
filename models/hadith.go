package models

// Hadith represents a single hadith with its number, Arabic text, and Indonesian translation
type Hadith struct {
	Number int    `json:"number"`
	Arab   string `json:"arab"`
	ID     string `json:"id"`
}

// HadithResponse is the standard response format for hadith API endpoints
type HadithResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is the standard error response format
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Pagination represents pagination information for list responses
type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
	PerPage     int `json:"per_page"`
}

// PaginatedResponse adds pagination information to the response
type PaginatedResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination Pagination  `json:"pagination"`
}

// Narrators contains the list of available narrators
type Narrators struct {
	Available []string `json:"available"`
}

// QueryParams represents the possible query parameters for filtering hadiths
type QueryParams struct {
	Page  int
	Limit int
	Query string
}