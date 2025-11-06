package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"tech-feed-weekly/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAllConfigs(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()

	// Create test config files
	testConfig1 := models.FeedData{
		Data: []models.FeedConfig{
			{
				Name:       "Test Feed 1",
				Type:       "categoryIsUrl",
				FeedURL:    "https://example1.com/feed.xml",
				LatestLink: "https://example1.com/latest1",
			},
		},
	}

	testConfig2 := models.FeedData{
		Data: []models.FeedConfig{
			{
				Name:       "Test Feed 2",
				Type:       "categoryIsAtomUrl",
				FeedURL:    "https://example2.com/feed.atom",
				LatestLink: "https://example2.com/latest2",
			},
			{
				Name:       "Test Feed 3",
				Type:       "hatena",
				FeedURL:    "https://example3.com",
				LatestLink: "https://example3.com/latest3",
			},
		},
	}

	// Write test config files
	config1Data, err := json.MarshalIndent(testConfig1, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tempDir, "test1.json"), config1Data, 0644)
	require.NoError(t, err)

	config2Data, err := json.MarshalIndent(testConfig2, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tempDir, "test2.json"), config2Data, 0644)
	require.NoError(t, err)

	// Create a non-JSON file (should be ignored)
	err = os.WriteFile(filepath.Join(tempDir, "readme.txt"), []byte("This is not a JSON file"), 0644)
	require.NoError(t, err)

	// Test LoadAllConfigs
	configMap, err := LoadAllConfigs(tempDir)
	require.NoError(t, err)

	// Verify results
	assert.Len(t, configMap, 2)

	// Check test1 config
	assert.Contains(t, configMap, "test1")
	test1Config := configMap["test1"]
	assert.Equal(t, filepath.Join(tempDir, "test1.json"), test1Config.FilePath)
	assert.Equal(t, "test1", test1Config.Category)
	assert.Len(t, test1Config.Data, 1)
	assert.Equal(t, "Test Feed 1", test1Config.Data[0].Name)
	assert.Equal(t, "test1", test1Config.Data[0].Category)

	// Check test2 config
	assert.Contains(t, configMap, "test2")
	test2Config := configMap["test2"]
	assert.Equal(t, filepath.Join(tempDir, "test2.json"), test2Config.FilePath)
	assert.Equal(t, "test2", test2Config.Category)
	assert.Len(t, test2Config.Data, 2)
	assert.Equal(t, "Test Feed 2", test2Config.Data[0].Name)
	assert.Equal(t, "Test Feed 3", test2Config.Data[1].Name)
	assert.Equal(t, "test2", test2Config.Data[0].Category)
	assert.Equal(t, "test2", test2Config.Data[1].Category)
}

func TestLoadAllConfigs_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	configMap, err := LoadAllConfigs(tempDir)
	require.NoError(t, err)
	assert.Empty(t, configMap)
}

func TestLoadAllConfigs_NonExistentDirectory(t *testing.T) {
	configMap, err := LoadAllConfigs("/non/existent/directory")
	assert.Error(t, err)
	assert.Nil(t, configMap)
}

func TestUpdateConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test.json")

	// Create initial config
	configData := &ConfigFileData{
		FilePath: configPath,
		Category: "test",
		Data: []models.FeedConfig{
			{
				Name:       "Test Feed",
				Type:       "categoryIsUrl",
				FeedURL:    "https://example.com/feed.xml",
				LatestLink: "https://example.com/old-latest",
				Category:   "test",
			},
		},
	}

	// Update the config
	configData.Data[0].LatestLink = "https://example.com/new-latest"

	// Test UpdateConfigFile
	err := UpdateConfigFile(configData)
	require.NoError(t, err)

	// Verify the file was updated
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var feedData models.FeedData
	err = json.Unmarshal(data, &feedData)
	require.NoError(t, err)

	assert.Len(t, feedData.Data, 1)
	assert.Equal(t, "https://example.com/new-latest", feedData.Data[0].LatestLink)
}

func TestGetAllFeedConfigs(t *testing.T) {
	configMap := map[string]*ConfigFileData{
		"test1": {
			FilePath: "/path/to/test1.json",
			Category: "test1",
			Data: []models.FeedConfig{
				{
					Name:     "Feed 1",
					Category: "test1",
				},
			},
		},
		"test2": {
			FilePath: "/path/to/test2.json",
			Category: "test2",
			Data: []models.FeedConfig{
				{
					Name:     "Feed 2",
					Category: "test2",
				},
				{
					Name:     "Feed 3",
					Category: "test2",
				},
			},
		},
	}

	allConfigs := GetAllFeedConfigs(configMap)

	assert.Len(t, allConfigs, 3)
	assert.Equal(t, "Feed 1", allConfigs[0].Name)
	assert.Equal(t, "test1", allConfigs[0].Category)

	// Note: order is not guaranteed in maps, so we check if all feeds are present
	feedNames := make([]string, len(allConfigs))
	for i, config := range allConfigs {
		feedNames[i] = config.Name
	}
	assert.ElementsMatch(t, []string{"Feed 1", "Feed 2", "Feed 3"}, feedNames)
}