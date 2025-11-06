package main

import (
	"log"
	"tech-feed-weekly/internal/config"
	"tech-feed-weekly/internal/feed"
	"tech-feed-weekly/internal/storage"
)

const (
	ConfigDir         = "config"
	LatestItemsPath   = "tmp/data/latest-items.json"
)

func main() {
	log.Println("Starting feed collector...")

	// Load all configuration files from configs directory
	log.Println("Loading configuration files...")
	configMap, err := config.LoadAllConfigs(ConfigDir)
	if err != nil {
		log.Fatalf("Failed to load configurations: %v", err)
	}

	totalFeeds := 0
	for _, configData := range configMap {
		totalFeeds += len(configData.Data)
	}
	log.Printf("Loaded %d configuration files with %d total feed configurations", len(configMap), totalFeeds)

	// Load existing latest items
	log.Println("Loading existing latest items...")
	existingItems, err := storage.LoadLatestItems(LatestItemsPath)
	if err != nil {
		log.Fatalf("Failed to load latest items: %v", err)
	}
	log.Printf("Loaded %d existing items", len(existingItems.Items))

	// Process all feeds to find new items and update config files
	log.Println("Processing feeds to find new items...")
	newItems, err := feed.ProcessAllFeeds(configMap, existingItems)
	if err != nil {
		log.Printf("Warning: Some feeds failed to process: %v", err)
		// Continue processing even if some feeds failed
	}

	if len(newItems) == 0 {
		log.Println("No new items found")
		return
	}

	log.Printf("Found %d new items", len(newItems))

	// Add new items to existing items
	itemsAdded := 0
	for _, newItem := range newItems {
		if storage.AddLatestItem(existingItems, newItem) {
			itemsAdded++
			log.Printf("Added new item: %s - %s", newItem.Category, newItem.Title)
		}
	}

	// Save updated items back to file
	if itemsAdded > 0 {
		log.Printf("Saving %d new items to %s", itemsAdded, LatestItemsPath)
		if err := storage.SaveLatestItems(LatestItemsPath, existingItems); err != nil {
			log.Fatalf("Failed to save latest items: %v", err)
		}
		log.Printf("Successfully saved %d new items", itemsAdded)
	} else {
		log.Println("No new items to save")
	}

	log.Println("Feed collector completed successfully")
}