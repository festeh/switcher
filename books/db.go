package books

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func GetDatabasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".local/share/zathura/bookmarks.sqlite"), nil
}

func GetLibraryDatabasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}
	
	// Create the directory if it doesn't exist
	dbDir := filepath.Join(homeDir, ".local/share/booklib")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", fmt.Errorf("error creating database directory: %w", err)
	}
	
	return filepath.Join(dbDir, "library.sqlite"), nil
}

func LoadDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	return db, nil
}

type BookmarkExtractor struct {
	DB *sql.DB
}

func NewBookmarkExtractor(path string) (*BookmarkExtractor, error) {
	db, err := LoadDatabase(path)
	if err != nil {
		return nil, err
	}
	return &BookmarkExtractor{DB: db}, nil
}

type BookmarkInfo struct {
	Filename string `json:"filename"`
	Page     int    `json:"page"`
	Title    string `json:"title"`
}

func getFileTitle(filePath string) string {
	cmd := exec.Command("exiftool", "-s", "-s", "-s", "-Title", filePath)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing exiftool on %s: %v", filePath, err)
		return filepath.Base(filePath) // Return filename as fallback
	}
	return string(output)
}

func (be *BookmarkExtractor) ExtractBookmarks() ([]BookmarkInfo, error) {
	rows, err := be.DB.Query("SELECT file, page FROM fileinfo")
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	var bookmarks []BookmarkInfo
	for rows.Next() {
		var filePath string
		var page int
		err := rows.Scan(&filePath, &page)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Printf("File does not exist: %s", filePath)
			continue
		}

		title := getFileTitle(filePath)

		bookmark := BookmarkInfo{
			Filename: filePath,
			Page:     page,
			Title:    title,
		}

		bookmarks = append(bookmarks, bookmark)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return bookmarks, nil
}
