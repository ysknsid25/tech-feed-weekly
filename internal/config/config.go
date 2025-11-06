package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"tech-feed-weekly/pkg/models"
)

// ConfigFileData represents config file data with filename
type ConfigFileData struct {
	FilePath string
	Category string
	Data     []models.FeedConfig
}

// LoadAllConfigs loads all JSON configuration files from the configs directory
// Returns a map with filename (without extension) as key and ConfigFileData as value
func LoadAllConfigs(configDir string) (map[string]*ConfigFileData, error) {
	configMap := make(map[string]*ConfigFileData)

	err := filepath.WalkDir(configDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-JSON files
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".json") {
			return nil
		}

		// Get filename without extension for category
		filename := d.Name()
		category := strings.TrimSuffix(filename, filepath.Ext(filename))

		// Read and parse JSON file
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read config file %s: %w", path, err)
		}

		var feedData models.FeedData
		if err := json.Unmarshal(data, &feedData); err != nil {
			return fmt.Errorf("failed to parse JSON file %s: %w", path, err)
		}

		// Set category for each feed config
		for i := range feedData.Data {
			feedData.Data[i].Category = category
		}

		configMap[category] = &ConfigFileData{
			FilePath: path,
			Category: category,
			Data:     feedData.Data,
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load configs: %w", err)
	}

	return configMap, nil
}

// UpdateConfigFile updates a specific config file with new latest links
func UpdateConfigFile(configData *ConfigFileData) error {
	feedData := models.FeedData{
		Data: configData.Data,
	}

	data, err := json.MarshalIndent(feedData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config data: %w", err)
	}

	if err := os.WriteFile(configData.FilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configData.FilePath, err)
	}

	return nil
}

// GetAllFeedConfigs extracts all feed configs from the config map
func GetAllFeedConfigs(configMap map[string]*ConfigFileData) []models.FeedConfig {
	var allConfigs []models.FeedConfig
	for _, configData := range configMap {
		allConfigs = append(allConfigs, configData.Data...)
	}
	return allConfigs
}