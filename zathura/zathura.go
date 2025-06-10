package zathura

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func NewZathura(path string) (*Zathura, error) {
	db, err := util.LoadDatabase(path)
	if err != nil {
		return nil, err
	}
	return &Zathura{DB: db}, nil
}

type BookInfo struct {
	Filename string `json:"filename"`
	Page     int    `json:"page"`
	Title    string `json:"title"`
}

func GetTitle(filePath string) string {
	cmd := exec.Command("exiftool", "-s", "-s", "-s", "-Title", filePath)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing exiftool on %s: %v", filePath, err)
		return filepath.Base(filePath) // Return filename as fallback
	}
	return string(output)
}

func (be *Zathura) GetAllKnownBooks() ([]BookInfo, error) {
	rows, err := be.DB.Query("SELECT file, page FROM fileinfo")
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	var bookmarks []BookInfo
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

		title := GetTitle(filePath)

		bookmark := BookInfo{
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
