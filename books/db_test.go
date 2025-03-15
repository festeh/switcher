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
