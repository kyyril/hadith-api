package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hadith-api/models"
)

// FileRepository handles loading and retrieving hadith data from JSON files
type FileRepository struct {
	DataDir string
	cache   map[string][]models.Hadith
}

// NewFileRepository creates a new file repository with the specified data directory
func NewFileRepository(dataDir string) *FileRepository {
	// Create data directory if it doesn't exist
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.MkdirAll(dataDir, 0755)
	}

	return &FileRepository{
		DataDir: dataDir,
		cache:   make(map[string][]models.Hadith),
	}
}

// GetAvailableNarrators returns a list of available narrators based on JSON files
func (r *FileRepository) GetAvailableNarrators() ([]string, error) {
	files, err := os.ReadDir(r.DataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	var narrators []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			narrator := strings.TrimSuffix(file.Name(), ".json")
			narrators = append(narrators, narrator)
		}
	}

	return narrators, nil
}

// GetHadithsByNarrator returns all hadiths from a specific narrator
func (r *FileRepository) GetHadithsByNarrator(narrator string, params models.QueryParams) ([]models.Hadith, int, error) {
	// Load data for the narrator
	hadiths, err := r.loadNarratorData(narrator)
	if err != nil {
		return nil, 0, err
	}

	// Apply filtering if query parameter is provided
	var filteredHadiths []models.Hadith
	if params.Query != "" {
		for _, h := range hadiths {
			if strings.Contains(strings.ToLower(h.ID), strings.ToLower(params.Query)) ||
				strings.Contains(strings.ToLower(h.Arab), strings.ToLower(params.Query)) {
				filteredHadiths = append(filteredHadiths, h)
			}
		}
	} else {
		filteredHadiths = hadiths
	}

	// Get total count before pagination
	totalItems := len(filteredHadiths)

	// Apply pagination if parameters are valid
	if params.Page > 0 && params.Limit > 0 {
		startIndex := (params.Page - 1) * params.Limit
		endIndex := startIndex + params.Limit

		// Check if startIndex is valid
		if startIndex >= totalItems {
			return []models.Hadith{}, totalItems, nil
		}

		// Check if endIndex is valid
		if endIndex > totalItems {
			endIndex = totalItems
		}

		filteredHadiths = filteredHadiths[startIndex:endIndex]
	}

	return filteredHadiths, totalItems, nil
}

// GetHadithByNumber returns a specific hadith by narrator and number
func (r *FileRepository) GetHadithByNumber(narrator string, number int) (*models.Hadith, error) {
	// Load data for the narrator
	hadiths, err := r.loadNarratorData(narrator)
	if err != nil {
		return nil, err
	}

	// Find hadith with the specified number
	for _, h := range hadiths {
		if h.Number == number {
			return &h, nil
		}
	}

	return nil, fmt.Errorf("hadith number %d not found for narrator %s", number, narrator)
}

// loadNarratorData loads hadith data for a specific narrator from the JSON file
func (r *FileRepository) loadNarratorData(narrator string) ([]models.Hadith, error) {
	// Check if the data is already cached
	if hadiths, ok := r.cache[narrator]; ok {
		return hadiths, nil
	}

	// Check if the file exists
	filePath := filepath.Join(r.DataDir, fmt.Sprintf("%s.json", narrator))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("narrator %s not found", narrator)
	}

	// Read the file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file for narrator %s: %w", narrator, err)
	}

	// Parse the JSON
	var hadiths []models.Hadith
	if err := json.Unmarshal(fileData, &hadiths); err != nil {
		return nil, fmt.Errorf("failed to parse JSON for narrator %s: %w", narrator, err)
	}

	// Cache the data
	r.cache[narrator] = hadiths

	return hadiths, nil
}