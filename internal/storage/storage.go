package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"tech-feed-weekly/pkg/models"
)

// LoadLatestItems loads existing latest items from the JSON file
func LoadLatestItems(filePath string) (*models.LatestItems, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create empty file with empty items array
		emptyItems := &models.LatestItems{Items: []models.LatestItem{}}
		if err := SaveLatestItems(filePath, emptyItems); err != nil {
			return nil, fmt.Errorf("failed to create new latest items file: %w", err)
		}
		return emptyItems, nil
	}

	// Read existing file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read latest items file: %w", err)
	}

	var latestItems models.LatestItems
	if err := json.Unmarshal(data, &latestItems); err != nil {
		return nil, fmt.Errorf("failed to parse latest items JSON: %w", err)
	}

	return &latestItems, nil
}

// SaveLatestItems saves latest items to the JSON file
func SaveLatestItems(filePath string, items *models.LatestItems) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal latest items: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write latest items file: %w", err)
	}

	return nil
}

// AddLatestItem adds a new item to latest items if it doesn't already exist
func AddLatestItem(items *models.LatestItems, newItem models.LatestItem) bool {
	// Check if item already exists
	for _, item := range items.Items {
		if item.Link == newItem.Link {
			return false // Item already exists
		}
	}

	// Add new item
	items.Items = append(items.Items, newItem)
	return true // Item was added
}