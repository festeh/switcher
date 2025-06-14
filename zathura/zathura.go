package zathura

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"switcher/util"
)

func GetDatabasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".local/share/zathura/bookmarks.sqlite"), nil
}

type Zathura struct {
	DB *sql.DB
}

func NewZathura() (*Zathura, error) {
	path, err := GetDatabasePath()
	if err != nil {
		return nil, err
	}
	db, err := util.LoadDatabase(path)
	if err != nil {
		return nil, err
	}
	return &Zathura{DB: db}, nil
}

type BookInfo struct {
	Filename string `json:"filename"`
	Page     int    `json:"page"`
}

func (zat *Zathura) GetAllKnownBooks() (map[string]BookInfo, error) {
	rows, err := zat.DB.Query("SELECT file, page FROM fileinfo")
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	bookmarks := make(map[string]BookInfo)
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

		bookmark := BookInfo{
			Filename: filePath,
			Page:     page,
		}

		bookmarks[filePath] = bookmark
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return bookmarks, nil
}
