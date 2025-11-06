package feed

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"tech-feed-weekly/pkg/models"
	"time"
)

// FetchLatestItem fetches the latest item from RSS/Atom feed
func FetchLatestItem(feedConfig models.FeedConfig) (*models.LatestItem, error) {
	feedURL := getFeedURL(feedConfig)
	if feedURL == "" {
		return nil, fmt.Errorf("could not generate feed URL for %s", feedConfig.Name)
	}

	isAtomFormat := isAtomFeed(feedConfig.Type)

	// Fetch the feed
	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed %s: %w", feedURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error %d when fetching %s", resp.StatusCode, feedURL)
	}

	if isAtomFormat {
		return parseAtomFeed(resp, feedConfig)
	}
	return parseRSSFeed(resp, feedConfig)
}

// getFeedURL generates the appropriate feed URL based on the feed type
func getFeedURL(config models.FeedConfig) string {
	switch config.Type {
	case "zenn":
		return fmt.Sprintf("https://zenn.dev/%s/feed", config.FeedURL)
	case "note":
		return fmt.Sprintf("https://note.com/%s/rss", config.FeedURL)
	case "qiita":
		return fmt.Sprintf("https://qiita.com/%s/feed", config.FeedURL)
	case "hatena":
		return fmt.Sprintf("%s/rss", config.FeedURL)
	case "scrapbox":
		return fmt.Sprintf("https://scrapbox.io/api/feed/%s", config.FeedURL)
	case "connpass":
		return fmt.Sprintf("https://%s.connpass.com/ja.atom", config.FeedURL)
	case "categoryIsUrl", "categoryIsAtomUrl":
		return config.FeedURL
	default:
		return ""
	}
}

// isAtomFeed determines if the feed type is Atom format
func isAtomFeed(feedType string) bool {
	return feedType == "qiita" || feedType == "connpass" || feedType == "categoryIsAtomUrl"
}

// parseRSSFeed parses RSS feed and returns the latest item
func parseRSSFeed(resp *http.Response, config models.FeedConfig) (*models.LatestItem, error) {
	var feed models.RSSFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("failed to decode RSS feed: %w", err)
	}

	if len(feed.Channel.Items) == 0 {
		return nil, fmt.Errorf("no items found in RSS feed")
	}

	// Parse dates and find the latest item
	var items []models.RSSItem
	for _, item := range feed.Channel.Items {
		parsedDate, err := parseDate(item.PubDate)
		if err != nil {
			// If date parsing fails, skip this item or use current time
			parsedDate = time.Now()
		}
		item.Date = parsedDate
		items = append(items, item)
	}

	// Sort by date (latest first)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Date.After(items[j].Date)
	})

	latestItem := items[0]
	return &models.LatestItem{
		Title:    strings.TrimSpace(latestItem.Title),
		Link:     strings.TrimSpace(latestItem.Link),
		Category: config.Category,
	}, nil
}

// parseAtomFeed parses Atom feed and returns the latest item
func parseAtomFeed(resp *http.Response, config models.FeedConfig) (*models.LatestItem, error) {
	var feed models.AtomFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("failed to decode Atom feed: %w", err)
	}

	if len(feed.Entries) == 0 {
		return nil, fmt.Errorf("no entries found in Atom feed")
	}

	// Parse dates and find the latest entry
	var entries []models.AtomEntry
	for _, entry := range feed.Entries {
		parsedDate, err := parseDate(entry.Updated)
		if err != nil {
			// If date parsing fails, skip this entry or use current time
			parsedDate = time.Now()
		}
		entry.Date = parsedDate
		entries = append(entries, entry)
	}

	// Sort by date (latest first)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Date.After(entries[j].Date)
	})

	latestEntry := entries[0]
	return &models.LatestItem{
		Title:    strings.TrimSpace(latestEntry.Title),
		Link:     strings.TrimSpace(latestEntry.Link.Href),
		Category: config.Category,
	}, nil
}

// parseDate parses various date formats commonly used in RSS/Atom feeds
func parseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)

	// Common date formats in RSS/Atom feeds
	formats := []string{
		time.RFC1123,     // RSS format: Mon, 02 Jan 2006 15:04:05 MST
		time.RFC1123Z,    // RSS format with numeric timezone
		time.RFC3339,     // Atom format: 2006-01-02T15:04:05Z07:00
		time.RFC3339Nano, // Atom format with nanoseconds
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 -0700",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// FetchHatenaBookmarkTechCategoryItems fetches items from Hatena Bookmark tech category RSS
// and filters them based on bookmark count and site-specific thresholds
func FetchHatenaBookmarkTechCategoryItems() ([]models.LatestItem, error) {
	url := "https://b.hatena.ne.jp/hotentry/it.rss"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Hatena Bookmark RSS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error %d when fetching Hatena Bookmark RSS", resp.StatusCode)
	}

	var feed models.HatenaBookmarkFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("failed to decode Hatena Bookmark RSS feed: %w", err)
	}

	if len(feed.Items) == 0 {
		return []models.LatestItem{}, nil
	}

	// Interest sites with lower threshold (matches TypeScript implementation)
	interestedSites := []string{
		"https://speakerdeck.com/",
	}

	var filteredItems []models.LatestItem
	for _, item := range feed.Items {
		// Filter out Zenn links to avoid duplication (matches TypeScript implementation)
		if strings.HasPrefix(item.Link, "https://zenn.dev/") {
			continue
		}

		// Check bookmark count with site-specific thresholds (matches TypeScript implementation)
		isInterestedSite := false
		for _, site := range interestedSites {
			if strings.HasPrefix(item.Link, site) {
				isInterestedSite = true
				break
			}
		}

		// Apply different thresholds based on site (matches TypeScript implementation)
		if isInterestedSite && item.BookmarkCount > 100 {
			filteredItems = append(filteredItems, models.LatestItem{
				Title:    strings.TrimSpace(item.Title),
				Link:     strings.TrimSpace(item.Link),
				Category: "hatena-bookmark-tech",
			})
		} else if !isInterestedSite && item.BookmarkCount > 120 {
			filteredItems = append(filteredItems, models.LatestItem{
				Title:    strings.TrimSpace(item.Title),
				Link:     strings.TrimSpace(item.Link),
				Category: "hatena-bookmark-tech",
			})
		}
	}

	return filteredItems, nil
}