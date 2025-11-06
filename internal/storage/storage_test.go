package storage

import (
	"os"
	"path/filepath"
	"testing"
	"tech-feed-weekly/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadLatestItems_NewFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "latest-items.json")

	// Test loading non-existent file (should create new one)
	items, err := LoadLatestItems(filePath)
	require.NoError(t, err)
	assert.NotNil(t, items)
	assert.Empty(t, items.Items)

	// Verify file was created
	_, err = os.Stat(filePath)
	assert.NoError(t, err)
}

func TestLoadLatestItems_ExistingFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "latest-items.json")

	// Create test data
	testItems := &models.LatestItems{
		Items: []models.LatestItem{
			{
				Title:    "Test Article 1",
				Link:     "https://example.com/article1",
				Category: "test1",
			},
			{
				Title:    "Test Article 2",
				Link:     "https://example.com/article2",
				Category: "test2",
			},
		},
	}

	// Save test data
	err := SaveLatestItems(filePath, testItems)
	require.NoError(t, err)

	// Load and verify
	loadedItems, err := LoadLatestItems(filePath)
	require.NoError(t, err)
	assert.Len(t, loadedItems.Items, 2)
	assert.Equal(t, "Test Article 1", loadedItems.Items[0].Title)
	assert.Equal(t, "https://example.com/article1", loadedItems.Items[0].Link)
	assert.Equal(t, "test1", loadedItems.Items[0].Category)
	assert.Equal(t, "Test Article 2", loadedItems.Items[1].Title)
	assert.Equal(t, "https://example.com/article2", loadedItems.Items[1].Link)
	assert.Equal(t, "test2", loadedItems.Items[1].Category)
}

func TestLoadLatestItems_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "invalid.json")

	// Create invalid JSON file
	err := os.WriteFile(filePath, []byte(`{"invalid": json`), 0644)
	require.NoError(t, err)

	// Test loading invalid JSON
	items, err := LoadLatestItems(filePath)
	assert.Error(t, err)
	assert.Nil(t, items)
}

func TestSaveLatestItems(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "save-test.json")

	testItems := &models.LatestItems{
		Items: []models.LatestItem{
			{
				Title:    "Save Test Article",
				Link:     "https://example.com/save-test",
				Category: "save-test",
			},
		},
	}

	// Test saving
	err := SaveLatestItems(filePath, testItems)
	require.NoError(t, err)

	// Verify file exists and has correct content
	_, err = os.Stat(filePath)
	assert.NoError(t, err)

	// Load and verify content
	loadedItems, err := LoadLatestItems(filePath)
	require.NoError(t, err)
	assert.Len(t, loadedItems.Items, 1)
	assert.Equal(t, "Save Test Article", loadedItems.Items[0].Title)
}

func TestSaveLatestItems_InvalidPath(t *testing.T) {
	// Test saving to invalid path
	testItems := &models.LatestItems{
		Items: []models.LatestItem{},
	}

	err := SaveLatestItems("/invalid/path/that/does/not/exist/file.json", testItems)
	assert.Error(t, err)
}

func TestAddLatestItem_NewItem(t *testing.T) {
	items := &models.LatestItems{
		Items: []models.LatestItem{
			{
				Title:    "Existing Article",
				Link:     "https://example.com/existing",
				Category: "existing",
			},
		},
	}

	newItem := models.LatestItem{
		Title:    "New Article",
		Link:     "https://example.com/new",
		Category: "new",
	}

	// Test adding new item
	added := AddLatestItem(items, newItem)
	assert.True(t, added)
	assert.Len(t, items.Items, 2)
	assert.Equal(t, "New Article", items.Items[1].Title)
	assert.Equal(t, "https://example.com/new", items.Items[1].Link)
	assert.Equal(t, "new", items.Items[1].Category)
}

func TestAddLatestItem_DuplicateItem(t *testing.T) {
	existingItem := models.LatestItem{
		Title:    "Existing Article",
		Link:     "https://example.com/existing",
		Category: "existing",
	}

	items := &models.LatestItems{
		Items: []models.LatestItem{existingItem},
	}

	// Test adding duplicate item (same link)
	duplicateItem := models.LatestItem{
		Title:    "Different Title",
		Link:     "https://example.com/existing", // Same link
		Category: "different",
	}

	added := AddLatestItem(items, duplicateItem)
	assert.False(t, added)
	assert.Len(t, items.Items, 1) // Should remain unchanged
	assert.Equal(t, "Existing Article", items.Items[0].Title) // Original item unchanged
}

func TestAddLatestItem_EmptyItems(t *testing.T) {
	items := &models.LatestItems{
		Items: []models.LatestItem{},
	}

	newItem := models.LatestItem{
		Title:    "First Article",
		Link:     "https://example.com/first",
		Category: "first",
	}

	// Test adding to empty items
	added := AddLatestItem(items, newItem)
	assert.True(t, added)
	assert.Len(t, items.Items, 1)
	assert.Equal(t, "First Article", items.Items[0].Title)
}