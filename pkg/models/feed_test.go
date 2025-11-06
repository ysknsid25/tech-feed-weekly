package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFeedConfig(t *testing.T) {
	feedConfig := FeedConfig{
		Name:            "Test Feed",
		Type:            "categoryIsUrl",
		FeedURL:         "https://example.com/feed.xml",
		LatestLink:      "https://example.com/latest",
		Category:        "test",
	}

	assert.Equal(t, "Test Feed", feedConfig.Name)
	assert.Equal(t, "categoryIsUrl", feedConfig.Type)
	assert.Equal(t, "https://example.com/feed.xml", feedConfig.FeedURL)
	assert.Equal(t, "https://example.com/latest", feedConfig.LatestLink)
	assert.Equal(t, "test", feedConfig.Category)
}

func TestFeedData(t *testing.T) {
	feedData := FeedData{
		Data: []FeedConfig{
			{
				Name:       "Feed 1",
				Type:       "categoryIsUrl",
				FeedURL:    "https://example1.com/feed.xml",
				LatestLink: "https://example1.com/latest",
				Category:   "test1",
			},
			{
				Name:       "Feed 2",
				Type:       "categoryIsAtomUrl",
				FeedURL:    "https://example2.com/feed.atom",
				LatestLink: "https://example2.com/latest",
				Category:   "test2",
			},
		},
	}

	assert.Len(t, feedData.Data, 2)
	assert.Equal(t, "Feed 1", feedData.Data[0].Name)
	assert.Equal(t, "Feed 2", feedData.Data[1].Name)
}

func TestLatestItem(t *testing.T) {
	item := LatestItem{
		Title:    "Test Article",
		Link:     "https://example.com/article",
		Category: "test",
	}

	assert.Equal(t, "Test Article", item.Title)
	assert.Equal(t, "https://example.com/article", item.Link)
	assert.Equal(t, "test", item.Category)
}

func TestLatestItems(t *testing.T) {
	items := LatestItems{
		Items: []LatestItem{
			{
				Title:    "Article 1",
				Link:     "https://example.com/article1",
				Category: "test1",
			},
			{
				Title:    "Article 2",
				Link:     "https://example.com/article2",
				Category: "test2",
			},
		},
	}

	assert.Len(t, items.Items, 2)
	assert.Equal(t, "Article 1", items.Items[0].Title)
	assert.Equal(t, "Article 2", items.Items[1].Title)
}

func TestRSSItem(t *testing.T) {
	testDate := time.Now()
	item := RSSItem{
		Title:   "RSS Article",
		Link:    "https://example.com/rss-article",
		PubDate: "Mon, 02 Jan 2006 15:04:05 MST",
		Date:    testDate,
	}

	assert.Equal(t, "RSS Article", item.Title)
	assert.Equal(t, "https://example.com/rss-article", item.Link)
	assert.Equal(t, "Mon, 02 Jan 2006 15:04:05 MST", item.PubDate)
	assert.Equal(t, testDate, item.Date)
}

func TestAtomEntry(t *testing.T) {
	testDate := time.Now()
	entry := AtomEntry{
		Title: "Atom Article",
		Link: AtomLink{
			Href: "https://example.com/atom-article",
		},
		Updated: "2006-01-02T15:04:05Z",
		Date:    testDate,
	}

	assert.Equal(t, "Atom Article", entry.Title)
	assert.Equal(t, "https://example.com/atom-article", entry.Link.Href)
	assert.Equal(t, "2006-01-02T15:04:05Z", entry.Updated)
	assert.Equal(t, testDate, entry.Date)
}