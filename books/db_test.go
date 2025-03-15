package books

import (
	"os"
	"testing"
)

func TestLoadDatabase(t *testing.T) {
	// Get the database path
	dbPath, err := GetDatabasePath()
	if err != nil {
		t.Fatalf("Failed to get database path: %v", err)
	}

	// Check if the database file exists
	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		t.Skipf("Database file does not exist at %s, skipping test", dbPath)
		return
	} else if err != nil {
		t.Fatalf("Error checking database file: %v", err)
	}

	// Try to load the database
	db, err := LoadDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to load database: %v", err)
	}
	defer db.Close()
}

func TestExtractBookmarks(t *testing.T) {
	// Get the database path
	dbPath, err := GetDatabasePath()
	if err != nil {
		t.Fatalf("Failed to get database path: %v", err)
	}

	// Check if the database file exists
	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		t.Skipf("Database file does not exist at %s, skipping test", dbPath)
		return
	} else if err != nil {
		t.Fatalf("Error checking database file: %v", err)
	}

	// Create a new bookmark extractor
	extractor, err := NewBookmarkExtractor(dbPath)
	if err != nil {
		t.Fatalf("Failed to create bookmark extractor: %v", err)
	}
	defer extractor.DB.Close()

	// Extract bookmarks
	bookmarks, err := extractor.ExtractBookmarks()
	if err != nil {
		t.Fatalf("Failed to extract bookmarks: %v", err)
	}

	// Print the results
	t.Logf("Successfully extracted %d bookmarks", len(bookmarks))
	for i, bookmark := range bookmarks {
		t.Logf("Bookmark %d: Title=%s, Filename=%s, Page=%d", 
			i+1, bookmark.Title, bookmark.Filename, bookmark.Page)
	}
}
