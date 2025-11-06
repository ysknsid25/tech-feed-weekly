package feed

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"tech-feed-weekly/internal/config"
	"tech-feed-weekly/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessFeedConfig_NewItem(t *testing.T) {
	// Mock RSS XML response
	rssXML := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <item>
      <title>New Article</title>
      <link>https://example.com/new-article</link>
      <pubDate>Mon, 06 Nov 2023 10:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rssXML))
	}))
	defer server.Close()

	// Test config with different latest link
	feedConfig := &models.FeedConfig{
		Name:       "Test Feed",
		Type:       "categoryIsUrl",
		FeedURL:    server.URL,
		LatestLink: "https://example.com/old-article", // Different from the new article
		Category:   "test",
	}

	existingItems := &models.LatestItems{
		Items: []models.LatestItem{},
	}

	// Test ProcessFeedConfig
	result := ProcessFeedConfig(feedConfig, existingItems)
	require.NotNil(t, result)
	assert.NoError(t, result.Error)
	assert.True(t, result.ConfigUpdated)
	assert.NotNil(t, result.NewItem)
	assert.Equal(t, "New Article", result.NewItem.Title)
	assert.Equal(t, "https://example.com/new-article", result.NewItem.Link)
	assert.Equal(t, "test", result.NewItem.Category)

	// Verify config was updated
	assert.Equal(t, "https://example.com/new-article", feedConfig.LatestLink)
}

func TestProcessFeedConfig_NoChange(t *testing.T) {
	// Mock RSS XML response
	rssXML := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <item>
      <title>Same Article</title>
      <link>https://example.com/same-article</link>
      <pubDate>Mon, 06 Nov 2023 10:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rssXML))
	}))
	defer server.Close()

	// Test config with same latest link
	feedConfig := &models.FeedConfig{
		Name:       "Test Feed",
		Type:       "categoryIsUrl",
		FeedURL:    server.URL,
		LatestLink: "https://example.com/same-article", // Same as the article in feed
		Category:   "test",
	}

	existingItems := &models.LatestItems{
		Items: []models.LatestItem{},
	}

	// Test ProcessFeedConfig
	result := ProcessFeedConfig(feedConfig, existingItems)
	require.NotNil(t, result)
	assert.NoError(t, result.Error)
	assert.False(t, result.ConfigUpdated)
	assert.Nil(t, result.NewItem)
}

func TestProcessFeedConfig_ItemAlreadyExists(t *testing.T) {
	// Mock RSS XML response
	rssXML := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <item>
      <title>Existing Article</title>
      <link>https://example.com/existing-article</link>
      <pubDate>Mon, 06 Nov 2023 10:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rssXML))
	}))
	defer server.Close()

	// Test config
	feedConfig := &models.FeedConfig{
		Name:       "Test Feed",
		Type:       "categoryIsUrl",
		FeedURL:    server.URL,
		LatestLink: "https://example.com/old-article", // Different from feed
		Category:   "test",
	}

	// Existing items contains the same item
	existingItems := &models.LatestItems{
		Items: []models.LatestItem{
			{
				Title:    "Existing Article",
				Link:     "https://example.com/existing-article", // Same as in feed
				Category: "test",
			},
		},
	}

	// Test ProcessFeedConfig
	result := ProcessFeedConfig(feedConfig, existingItems)
	require.NotNil(t, result)
	assert.NoError(t, result.Error)
	assert.False(t, result.ConfigUpdated)
	assert.Nil(t, result.NewItem)
}

func TestProcessFeedConfig_FetchError(t *testing.T) {
	// Test config with invalid URL
	feedConfig := &models.FeedConfig{
		Name:       "Test Feed",
		Type:       "categoryIsUrl",
		FeedURL:    "http://invalid-url-that-does-not-exist.example",
		LatestLink: "https://example.com/old-article",
		Category:   "test",
	}

	existingItems := &models.LatestItems{
		Items: []models.LatestItem{},
	}

	// Test ProcessFeedConfig
	result := ProcessFeedConfig(feedConfig, existingItems)
	require.NotNil(t, result)
	assert.Error(t, result.Error)
	assert.False(t, result.ConfigUpdated)
	assert.Nil(t, result.NewItem)
}

func TestProcessAllFeeds(t *testing.T) {
	tempDir := t.TempDir()

	// Mock RSS XML responses
	rssXML1 := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <item>
      <title>New Article 1</title>
      <link>https://example1.com/new-article</link>
      <pubDate>Mon, 06 Nov 2023 10:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>`

	rssXML2 := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <item>
      <title>New Article 2</title>
      <link>https://example2.com/new-article</link>
      <pubDate>Mon, 06 Nov 2023 11:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>`

	// Create test servers
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rssXML1))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rssXML2))
	}))
	defer server2.Close()

	// Create config map
	configMap := map[string]*config.ConfigFileData{
		"test1": {
			FilePath: filepath.Join(tempDir, "test1.json"),
			Category: "test1",
			Data: []models.FeedConfig{
				{
					Name:       "Test Feed 1",
					Type:       "categoryIsUrl",
					FeedURL:    server1.URL,
					LatestLink: "https://example1.com/old-article",
					Category:   "test1",
				},
			},
		},
		"test2": {
			FilePath: filepath.Join(tempDir, "test2.json"),
			Category: "test2",
			Data: []models.FeedConfig{
				{
					Name:       "Test Feed 2",
					Type:       "categoryIsUrl",
					FeedURL:    server2.URL,
					LatestLink: "https://example2.com/old-article",
					Category:   "test2",
				},
			},
		},
	}

	existingItems := &models.LatestItems{
		Items: []models.LatestItem{},
	}

	// Test ProcessAllFeeds
	newItems, err := ProcessAllFeeds(configMap, existingItems)
	require.NoError(t, err)
	assert.Len(t, newItems, 2)

	// Verify new items
	assert.Equal(t, "New Article 1", newItems[0].Title)
	assert.Equal(t, "https://example1.com/new-article", newItems[0].Link)
	assert.Equal(t, "test1", newItems[0].Category)

	assert.Equal(t, "New Article 2", newItems[1].Title)
	assert.Equal(t, "https://example2.com/new-article", newItems[1].Link)
	assert.Equal(t, "test2", newItems[1].Category)

	// Verify configs were updated
	assert.Equal(t, "https://example1.com/new-article", configMap["test1"].Data[0].LatestLink)
	assert.Equal(t, "https://example2.com/new-article", configMap["test2"].Data[0].LatestLink)

	// Verify config files were created (UpdateConfigFile was called)
	_, err = os.Stat(configMap["test1"].FilePath)
	assert.NoError(t, err)
	_, err = os.Stat(configMap["test2"].FilePath)
	assert.NoError(t, err)
}

func TestProcessAllFeeds_NoNewItems(t *testing.T) {
	// Mock RSS XML response
	rssXML := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <item>
      <title>Same Article</title>
      <link>https://example.com/same-article</link>
      <pubDate>Mon, 06 Nov 2023 10:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rssXML))
	}))
	defer server.Close()

	configMap := map[string]*config.ConfigFileData{
		"test": {
			FilePath: "/tmp/test.json",
			Category: "test",
			Data: []models.FeedConfig{
				{
					Name:       "Test Feed",
					Type:       "categoryIsUrl",
					FeedURL:    server.URL,
					LatestLink: "https://example.com/same-article", // Same as in feed
					Category:   "test",
				},
			},
		},
	}

	existingItems := &models.LatestItems{
		Items: []models.LatestItem{},
	}

	// Test ProcessAllFeeds
	newItems, err := ProcessAllFeeds(configMap, existingItems)
	require.NoError(t, err)
	assert.Empty(t, newItems)
}

func TestProcessAllFeeds_WithErrors(t *testing.T) {
	tempDir := t.TempDir()

	configMap := map[string]*config.ConfigFileData{
		"test": {
			FilePath: filepath.Join(tempDir, "test.json"),
			Category: "test",
			Data: []models.FeedConfig{
				{
					Name:       "Invalid Feed",
					Type:       "categoryIsUrl",
					FeedURL:    "http://invalid-url.example",
					LatestLink: "https://example.com/old-article",
					Category:   "test",
				},
			},
		},
	}

	existingItems := &models.LatestItems{
		Items: []models.LatestItem{},
	}

	// Test ProcessAllFeeds with error
	newItems, err := ProcessAllFeeds(configMap, existingItems)
	assert.Error(t, err)
	assert.Empty(t, newItems)
}