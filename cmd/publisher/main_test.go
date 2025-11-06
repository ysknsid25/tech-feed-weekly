package main

import (
	"os"
	"path/filepath"
	"strings"
	"tech-feed-weekly/pkg/models"
	"testing"
)

func TestFormatCategoryName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hatena-bookmark-tech", "はてなブックマーク テック"},
		{"tech-articles", "Tech Articles"},
		{"golang", "Golang"},
		{"machine-learning-ai", "Machine Learning Ai"},
	}

	for _, test := range tests {
		result := formatCategoryName(test.input)
		if result != test.expected {
			t.Errorf("formatCategoryName(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "Hello World"},
		{"<script>alert('test')</script>", "&lt;script&gt;alert(&#39;test&#39;)&lt;/script&gt;"},
		{"A & B", "A &amp; B"},
		{`"quoted"`, "&quot;quoted&quot;"},
	}

	for _, test := range tests {
		result := escapeHTML(test.input)
		if result != test.expected {
			t.Errorf("escapeHTML(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestGenerateHTML(t *testing.T) {
	// Create test data
	latestItems := &models.LatestItems{
		Items: []models.LatestItem{
			{
				Title:    "Test Article 1",
				Link:     "https://example.com/article1",
				Category: "tech-articles",
			},
			{
				Title:    "Test Article 2",
				Link:     "https://example.com/article2",
				Category: "hatena-bookmark-tech",
			},
		},
	}

	// Generate HTML
	html := generateHTML(latestItems)

	// Verify HTML structure
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("Generated HTML should contain DOCTYPE declaration")
	}

	if !strings.Contains(html, "<h1>Tech Feed Weekly") {
		t.Error("Generated HTML should contain main title")
	}

	if !strings.Contains(html, "<h2>Tech Articles</h2>") {
		t.Error("Generated HTML should contain tech-articles section")
	}

	if !strings.Contains(html, "<h2>はてなブックマーク テック</h2>") {
		t.Error("Generated HTML should contain hatena-bookmark-tech section")
	}

	if !strings.Contains(html, `<a href="https://example.com/article1">Test Article 1</a>`) {
		t.Error("Generated HTML should contain article 1 link")
	}

	if !strings.Contains(html, `<a href="https://example.com/article2">Test Article 2</a>`) {
		t.Error("Generated HTML should contain article 2 link")
	}

	if !strings.Contains(html, "Total items: 2") {
		t.Error("Generated HTML should contain total items count")
	}
}

func TestLoadLatestItems(t *testing.T) {
	// Create temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-items.json")

	testData := `{
		"items": [
			{
				"title": "Test Article",
				"link": "https://example.com/test",
				"category": "test-category"
			}
		]
	}`

	err := os.WriteFile(testFile, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test loading
	items, err := loadLatestItems(testFile)
	if err != nil {
		t.Fatalf("Failed to load latest items: %v", err)
	}

	if len(items.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items.Items))
	}

	if items.Items[0].Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got '%s'", items.Items[0].Title)
	}

	if items.Items[0].Link != "https://example.com/test" {
		t.Errorf("Expected link 'https://example.com/test', got '%s'", items.Items[0].Link)
	}

	if items.Items[0].Category != "test-category" {
		t.Errorf("Expected category 'test-category', got '%s'", items.Items[0].Category)
	}
}

func TestLoadLatestItems_NonExistentFile(t *testing.T) {
	_, err := loadLatestItems("nonexistent-file.json")
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}
}

func TestLoadLatestItems_InvalidJSON(t *testing.T) {
	// Create temporary test file with invalid JSON
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "invalid.json")

	err := os.WriteFile(testFile, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = loadLatestItems(testFile)
	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
}