package models

import "time"

// FeedConfig represents a feed configuration loaded from JSON file
type FeedConfig struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	FeedURL         string `json:"feedUrl"`
	LatestLink      string `json:"latestLink"`
	Category        string `json:"-"` // File name without extension
}

// FeedData represents the data structure for feeds configuration
type FeedData struct {
	Data []FeedConfig `json:"data"`
}

// LatestItem represents an item in latest-items.json
type LatestItem struct {
	Title    string `json:"title"`
	Link     string `json:"link"`
	Category string `json:"category"`
}

// LatestItems represents the structure of latest-items.json
type LatestItems struct {
	Items []LatestItem `json:"items"`
}

// RSSItem represents an item from RSS feed
type RSSItem struct {
	Title   string    `xml:"title"`
	Link    string    `xml:"link"`
	PubDate string    `xml:"pubDate"`
	Date    time.Time `xml:"-"`
}

// AtomEntry represents an entry from Atom feed
type AtomEntry struct {
	Title   string    `xml:"title"`
	Link    AtomLink  `xml:"link"`
	Updated string    `xml:"updated"`
	Date    time.Time `xml:"-"`
}

// AtomLink represents a link in Atom feed
type AtomLink struct {
	Href string `xml:"href,attr"`
}

// RSSFeed represents RSS feed structure
type RSSFeed struct {
	Channel RSSChannel `xml:"channel"`
}

// RSSChannel represents RSS channel
type RSSChannel struct {
	Items []RSSItem `xml:"item"`
}

// AtomFeed represents Atom feed structure
type AtomFeed struct {
	Entries []AtomEntry `xml:"entry"`
}