package feed

import (
	"fmt"
	"log"
	"tech-feed-weekly/internal/config"
	"tech-feed-weekly/pkg/models"
)

// ProcessResult represents the result of processing a feed
type ProcessResult struct {
	NewItem       *models.LatestItem
	ConfigUpdated bool
	Error         error
}

// ProcessFeedConfig processes a single feed configuration and returns the latest item if it's new
// Also updates the config if a new item is found
func ProcessFeedConfig(config *models.FeedConfig, existingItems *models.LatestItems) *ProcessResult {
	result := &ProcessResult{}

	// Fetch the latest item from the feed
	latestItem, err := FetchLatestItem(*config)
	if err != nil {
		result.Error = fmt.Errorf("failed to fetch latest item for %s: %w", config.Name, err)
		return result
	}

	// Check if the latest item URL is different from the recorded latestLink
	if latestItem.Link == config.LatestLink {
		log.Printf("No new item for %s: latest link unchanged", config.Name)
		return result // No new item
	}

	// Check if this item already exists in the existing items
	for _, item := range existingItems.Items {
		if item.Link == latestItem.Link {
			log.Printf("Item already exists for %s: %s", config.Name, latestItem.Link)
			return result // Item already exists
		}
	}

	log.Printf("New item found for %s: %s", config.Name, latestItem.Title)

	// Update the config's LatestLink
	config.LatestLink = latestItem.Link
	result.NewItem = latestItem
	result.ConfigUpdated = true

	return result
}

// ProcessAllFeeds processes all feed configurations and returns new items
// Also updates config files when new items are found
func ProcessAllFeeds(configMap map[string]*config.ConfigFileData, existingItems *models.LatestItems) ([]models.LatestItem, error) {
	var newItems []models.LatestItem
	var errors []error
	updatedConfigs := make(map[string]bool)

	for categoryName, configData := range configMap {
		log.Printf("Processing category: %s", categoryName)

		for i := range configData.Data {
			feedConfig := &configData.Data[i]
			log.Printf("Processing feed: %s", feedConfig.Name)

			result := ProcessFeedConfig(feedConfig, existingItems)
			if result.Error != nil {
				log.Printf("Error processing %s: %v", feedConfig.Name, result.Error)
				errors = append(errors, result.Error)
				continue
			}

			if result.NewItem != nil {
				newItems = append(newItems, *result.NewItem)
			}

			if result.ConfigUpdated {
				updatedConfigs[categoryName] = true
			}
		}
	}

	// Update config files that had changes
	for categoryName := range updatedConfigs {
		if err := config.UpdateConfigFile(configMap[categoryName]); err != nil {
			log.Printf("Error updating config file for %s: %v", categoryName, err)
			errors = append(errors, fmt.Errorf("failed to update config file for %s: %w", categoryName, err))
		} else {
			log.Printf("Updated config file for category: %s", categoryName)
		}
	}

	if len(errors) > 0 {
		log.Printf("Encountered %d errors during processing", len(errors))
		// Return the first error for simplicity, but log all errors
		return newItems, fmt.Errorf("encountered errors during processing: %v", errors[0])
	}

	return newItems, nil
}