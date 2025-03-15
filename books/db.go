package books

import (
	"fmt"
	"os"
	"path/filepath"
  _ "github.com/mattn/go-sqlite3"
  "database/sql"
)

func GetDatabasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".local/share/zathura/bookmarks.sqlite"), nil
}

func LoadDatabase(path string) (*sql.DB, error) {
  db, err := sql.Open("sqlite3", path)
  if err != nil {
    return nil, fmt.Errorf("error opening database: %w", err)
  }
  return db, nil
}
