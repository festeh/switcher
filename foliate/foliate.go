package foliate

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetDataPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".local/share/com.github.johnfactotum.Foliate"), nil
}

type Foliate struct {
	DataPath string
}

func NewFoliate() (*Foliate, error) {
	path, err := GetDataPath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("foliate data directory does not exist: %s", path)
	}
	return &Foliate{DataPath: path}, nil
}

type BookInfo struct {
	Filename string `json:"filename"`
	Page     int    `json:"page"`
	Title    string `json:"title"`
	Author   string `json:"author,omitempty"`
}

type FoliateMetadata struct {
	Title      string `json:"title"`
	Identifier string `json:"identifier"`
	Language   string `json:"language"`
	Author     []struct {
		Name   string `json:"name"`
		SortAs string `json:"sortAs"`
	} `json:"author"`
	Publisher   string   `json:"publisher"`
	Description string   `json:"description"`
	Subject     []string `json:"subject"`
}

type FoliateBook struct {
	Metadata     FoliateMetadata `json:"metadata"`
	Progress     []int           `json:"progress"`
	LastLocation string          `json:"lastLocation"`
}

type URIStore struct {
	URIs [][]string `json:"uris"`
}

func (f *Foliate) GetAllKnownBooks() (map[string]BookInfo, error) {
	uriStorePath := filepath.Join(f.DataPath, "library", "uri-store.json")

	uriData, err := os.ReadFile(uriStorePath)
	if err != nil {
		return nil, fmt.Errorf("error reading uri-store.json: %w", err)
	}

	var uriStore URIStore
	if err := json.Unmarshal(uriData, &uriStore); err != nil {
		return nil, fmt.Errorf("error parsing uri-store.json: %w", err)
	}

	books := make(map[string]BookInfo)
	for _, uri := range uriStore.URIs {
		if len(uri) != 2 {
			continue
		}

		identifier := uri[0]
		filePath := uri[1]

		if strings.HasPrefix(filePath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Printf("Error getting home directory: %v", err)
				continue
			}
			filePath = filepath.Join(homeDir, filePath[2:])
		}

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Printf("File does not exist: %s", filePath)
			continue
		}

		bookMetadataPath := filepath.Join(f.DataPath, identifier+".json")
		bookData, err := os.ReadFile(bookMetadataPath)
		if err != nil {
			log.Printf("Error reading book metadata %s: %v", bookMetadataPath, err)
			continue
		}

		var foliateBook FoliateBook
		if err := json.Unmarshal(bookData, &foliateBook); err != nil {
			log.Printf("Error parsing book metadata %s: %v", bookMetadataPath, err)
			continue
		}

		title := foliateBook.Metadata.Title
		if title == "" {
			title = filepath.Base(filePath)
		}

		author := ""
		if len(foliateBook.Metadata.Author) > 0 {
			author = foliateBook.Metadata.Author[0].Name
		}

		page := 0
		if len(foliateBook.Progress) > 0 {
			page = foliateBook.Progress[0]
		}

		book := BookInfo{
			Filename: filePath,
			Page:     page,
			Title:    title,
			Author:   author,
		}

		books[filePath] = book
	}

	return books, nil
}
