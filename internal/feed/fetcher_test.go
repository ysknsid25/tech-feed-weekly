package feed

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tech-feed-weekly/pkg/models"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFeedURL(t *testing.T) {
	tests := []struct {
		name     string
		config   models.FeedConfig
		expected string
	}{
		{
			name: "zenn feed",
			config: models.FeedConfig{
				Type:    "zenn",
				FeedURL: "username",
			},
			expected: "https://zenn.dev/username/feed",
		},
		{
			name: "note feed",
			config: models.FeedConfig{
				Type:    "note",
				FeedURL: "username",
			},
			expected: "https://note.com/username/rss",
		},
		{
			name: "qiita feed",
			config: models.FeedConfig{
				Type:    "qiita",
				FeedURL: "username",
			},
			expected: "https://qiita.com/username/feed",
		},
		{
			name: "hatena feed",
			config: models.FeedConfig{
				Type:    "hatena",
				FeedURL: "https://example.hatenablog.com",
			},
			expected: "https://example.hatenablog.com/rss",
		},
		{
			name: "scrapbox feed",
			config: models.FeedConfig{
				Type:    "scrapbox",
				FeedURL: "projectname",
			},
			expected: "https://scrapbox.io/api/feed/projectname",
		},
		{
			name: "connpass feed",
			config: models.FeedConfig{
				Type:    "connpass",
				FeedURL: "groupname",
			},
			expected: "https://groupname.connpass.com/ja.atom",
		},
		{
			name: "categoryIsUrl",
			config: models.FeedConfig{
				Type:    "categoryIsUrl",
				FeedURL: "https://example.com/feed.xml",
			},
			expected: "https://example.com/feed.xml",
		},
		{
			name: "categoryIsAtomUrl",
			config: models.FeedConfig{
				Type:    "categoryIsAtomUrl",
				FeedURL: "https://example.com/feed.atom",
			},
			expected: "https://example.com/feed.atom",
		},
		{
			name: "unknown type",
			config: models.FeedConfig{
				Type:    "unknown",
				FeedURL: "test",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFeedURL(tt.config)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsAtomFeed(t *testing.T) {
	tests := []struct {
		feedType string
		expected bool
	}{
		{"qiita", true},
		{"connpass", true},
		{"categoryIsAtomUrl", true},
		{"zenn", false},
		{"note", false},
		{"hatena", false},
		{"scrapbox", false},
		{"categoryIsUrl", false},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.feedType, func(t *testing.T) {
			result := isAtomFeed(tt.feedType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name     string
		dateStr  string
		hasError bool
	}{
		{
			name:     "RFC1123",
			dateStr:  "Mon, 02 Jan 2006 15:04:05 MST",
			hasError: false,
		},
		{
			name:     "RFC1123Z",
			dateStr:  "Mon, 02 Jan 2006 15:04:05 -0700",
			hasError: false,
		},
		{
			name:     "RFC3339",
			dateStr:  "2006-01-02T15:04:05+07:00",
			hasError: false,
		},
		{
			name:     "RFC3339Nano",
			dateStr:  "2006-01-02T15:04:05.999999999+07:00",
			hasError: false,
		},
		{
			name:     "Simple ISO format",
			dateStr:  "2006-01-02T15:04:05Z",
			hasError: false,
		},
		{
			name:     "Simple ISO without Z",
			dateStr:  "2006-01-02T15:04:05",
			hasError: false,
		},
		{
			name:     "Simple datetime",
			dateStr:  "2006-01-02 15:04:05",
			hasError: false,
		},
		{
			name:     "RFC1123 variant",
			dateStr:  "Mon, 2 Jan 2006 15:04:05 MST",
			hasError: false,
		},
		{
			name:     "RFC1123 with timezone",
			dateStr:  "Mon, 2 Jan 2006 15:04:05 -0700",
			hasError: false,
		},
		{
			name:     "Invalid date",
			dateStr:  "invalid date format",
			hasError: true,
		},
		{
			name:     "Empty string",
			dateStr:  "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDate(tt.dateStr)
			if tt.hasError {
				assert.Error(t, err)
				assert.Equal(t, time.Time{}, result)
			} else {
				assert.NoError(t, err)
				assert.False(t, result.IsZero())
			}
		})
	}
}

func TestFetchLatestItem_RSS(t *testing.T) {
	// Mock RSS XML response
	rssXML := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test RSS Feed</title>
    <item>
      <title>Latest Article</title>
      <link>https://example.com/latest</link>
      <pubDate>Mon, 06 Nov 2023 10:00:00 GMT</pubDate>
    </item>
    <item>
      <title>Older Article</title>
      <link>https://example.com/older</link>
      <pubDate>Mon, 05 Nov 2023 10:00:00 GMT</pubDate>
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
	config := models.FeedConfig{
		Name:     "Test RSS Feed",
		Type:     "categoryIsUrl",
		FeedURL:  server.URL,
		Category: "test",
	}

	// Test FetchLatestItem
	item, err := FetchLatestItem(config)
	require.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "Latest Article", item.Title)
	assert.Equal(t, "https://example.com/latest", item.Link)
	assert.Equal(t, "test", item.Category)
}

func TestFetchLatestItem_Atom(t *testing.T) {
	// Mock Atom XML response
	atomXML := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Test Atom Feed</title>
  <entry>
    <title>Latest Entry</title>
    <link href="https://example.com/latest-entry"/>
    <updated>2023-11-06T10:00:00Z</updated>
  </entry>
  <entry>
    <title>Older Entry</title>
    <link href="https://example.com/older-entry"/>
    <updated>2023-11-05T10:00:00Z</updated>
  </entry>
</feed>`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/atom+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(atomXML))
	}))
	defer server.Close()

	// Test config
	config := models.FeedConfig{
		Name:     "Test Atom Feed",
		Type:     "categoryIsAtomUrl",
		FeedURL:  server.URL,
		Category: "test",
	}

	// Test FetchLatestItem
	item, err := FetchLatestItem(config)
	require.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "Latest Entry", item.Title)
	assert.Equal(t, "https://example.com/latest-entry", item.Link)
	assert.Equal(t, "test", item.Category)
}

func TestFetchLatestItem_HTTPError(t *testing.T) {
	// Create test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := models.FeedConfig{
		Name:     "Test Feed",
		Type:     "categoryIsUrl",
		FeedURL:  server.URL,
		Category: "test",
	}

	// Test FetchLatestItem
	item, err := FetchLatestItem(config)
	assert.Error(t, err)
	assert.Nil(t, item)
	assert.Contains(t, err.Error(), "HTTP error 404")
}

func TestFetchLatestItem_InvalidXML(t *testing.T) {
	// Create test server with invalid XML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid xml content"))
	}))
	defer server.Close()

	config := models.FeedConfig{
		Name:     "Test Feed",
		Type:     "categoryIsUrl",
		FeedURL:  server.URL,
		Category: "test",
	}

	// Test FetchLatestItem
	item, err := FetchLatestItem(config)
	assert.Error(t, err)
	assert.Nil(t, item)
	assert.Contains(t, strings.ToLower(err.Error()), "decode")
}

func TestFetchLatestItem_EmptyFeed(t *testing.T) {
	// Mock empty RSS XML response
	emptyRSSXML := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Empty RSS Feed</title>
  </channel>
</rss>`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(emptyRSSXML))
	}))
	defer server.Close()

	config := models.FeedConfig{
		Name:     "Empty Feed",
		Type:     "categoryIsUrl",
		FeedURL:  server.URL,
		Category: "test",
	}

	// Test FetchLatestItem
	item, err := FetchLatestItem(config)
	assert.Error(t, err)
	assert.Nil(t, item)
	assert.Contains(t, err.Error(), "no items found")
}