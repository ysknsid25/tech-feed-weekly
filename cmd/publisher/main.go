package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"tech-feed-weekly/pkg/models"
	"time"
)

const (
	LatestItemsPath = "tmp/data/latest-items.json"
	OutputDir       = "tmp/publisher"
	OutputFile      = "newsletter.html"
)

func main() {
	log.Println("Starting publisher...")

	// Create output directory
	if err := os.MkdirAll(OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Load latest items
	log.Println("Loading latest items...")
	latestItems, err := loadLatestItems(LatestItemsPath)
	if err != nil {
		log.Fatalf("Failed to load latest items: %v", err)
	}

	if len(latestItems.Items) == 0 {
		log.Println("No items found to publish")
		return
	}

	log.Printf("Found %d items to publish", len(latestItems.Items))

	// Generate HTML content
	log.Println("Generating HTML content...")
	htmlContent := generateHTML(latestItems)

	// Save to output file
	outputPath := filepath.Join(OutputDir, OutputFile)
	log.Printf("Saving content to %s", outputPath)
	if err := os.WriteFile(outputPath, []byte(htmlContent), 0644); err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	// Clean up latest items file
	log.Printf("Cleaning up %s", LatestItemsPath)
	if err := os.Remove(LatestItemsPath); err != nil {
		log.Printf("Warning: Failed to remove latest items file: %v", err)
	}

	log.Println("Publisher completed successfully")
}

// loadLatestItems loads the latest items from JSON file
func loadLatestItems(filePath string) (*models.LatestItems, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read latest items file: %w", err)
	}

	var latestItems models.LatestItems
	if err := json.Unmarshal(data, &latestItems); err != nil {
		return nil, fmt.Errorf("failed to parse latest items JSON: %w", err)
	}

	return &latestItems, nil
}

// generateHTML creates HTML content grouped by category
func generateHTML(latestItems *models.LatestItems) string {
	// Group items by category
	categoryMap := make(map[string][]models.LatestItem)
	for _, item := range latestItems.Items {
		categoryMap[item.Category] = append(categoryMap[item.Category], item)
	}

	// Sort categories for consistent output
	var categories []string
	for category := range categoryMap {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	var htmlBuilder strings.Builder

	// HTML header
	htmlBuilder.WriteString(`<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tech Feed Weekly</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; line-height: 1.6; max-width: 800px; margin: 0 auto; padding: 20px; }
        h1 { color: #333; border-bottom: 2px solid #007acc; padding-bottom: 10px; }
        h2 { color: #555; margin-top: 30px; }
        ul { list-style-type: none; padding: 0; }
        li { margin: 10px 0; padding: 8px; background-color: #f8f9fa; border-radius: 4px; }
        a { color: #007acc; text-decoration: none; }
        a:hover { text-decoration: underline; }
        .footer { margin-top: 40px; padding-top: 20px; border-top: 1px solid #ddd; color: #666; font-size: 0.9em; }
    </style>
</head>
<body>
    <h1>Tech Feed Weekly - ` + time.Now().Format("2006-01-02") + `</h1>
`)

	// Generate content for each category
	for _, category := range categories {
		items := categoryMap[category]

		// Category header
		htmlBuilder.WriteString(fmt.Sprintf("    <h2>%s</h2>\n", formatCategoryName(category)))
		htmlBuilder.WriteString("    <ul>\n")

		// Sort items by title for consistent output
		sort.Slice(items, func(i, j int) bool {
			return items[i].Title < items[j].Title
		})

		// Add items
		for _, item := range items {
			htmlBuilder.WriteString(fmt.Sprintf("        <li><a href=\"%s\">%s</a></li>\n",
				escapeHTML(item.Link), escapeHTML(item.Title)))
		}

		htmlBuilder.WriteString("    </ul>\n")
	}

	// HTML footer
	htmlBuilder.WriteString(`    <div class="footer">
        <p>Generated on ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
        <p>Total items: ` + fmt.Sprintf("%d", len(latestItems.Items)) + `</p>
    </div>
</body>
</html>`)

	return htmlBuilder.String()
}

// formatCategoryName formats category name for display
func formatCategoryName(category string) string {
	// Convert category names to more readable format
	switch category {
	case "hatena-bookmark-tech":
		return "はてなブックマーク テック"
	default:
		// Capitalize first letter and replace hyphens with spaces
		parts := strings.Split(category, "-")
		for i, part := range parts {
			if len(part) > 0 {
				parts[i] = strings.ToUpper(part[:1]) + part[1:]
			}
		}
		return strings.Join(parts, " ")
	}
}

// escapeHTML escapes HTML special characters
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}