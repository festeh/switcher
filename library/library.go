package library

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"switcher/util"
	"switcher/zathura"
)

type Book struct {
	ID       int64     `json:"id"`
	FilePath string    `json:"filepath"`
	Title    string    `json:"title"`
	FileSize int64     `json:"filesize"`
	ModTime  time.Time `json:"modtime"`
	Format   string    `json:"format"`
}

type Library struct {
	DB        *sql.DB
	Extractor *zathura.BookmarkExtractor
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

func NewLibrary(dbPath string) (*Library, error) {
	db, err := util.LoadDatabase(dbPath)
	if err != nil {
		return nil, err
	}

	library := &Library{DB: db}
	if err := library.initSchema(); err != nil {
		return nil, err
	}

	// Initialize the Zathura bookmark extractor
	if err := library.initBookmarkExtractor(); err != nil {
		fmt.Printf("Failed to initialize bookmark extractor: %v\n", err)
		// Continue without extractor - it's not critical for library functionality
	}

	return library, nil
}

func (l *Library) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filepath TEXT UNIQUE NOT NULL,
		title TEXT NOT NULL,
		filesize INTEGER NOT NULL,
		modtime DATETIME NOT NULL,
		format TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_books_filepath ON books(filepath);
	CREATE INDEX IF NOT EXISTS idx_books_format ON books(format);
	`
	_, err := l.DB.Exec(query)
	return err
}

func (l *Library) ScanDirectory(rootDir string) error {
	supportedFormats := map[string]bool{
		".pdf":  true,
		".epub": true,
		".fb2":  true,
	}

	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !supportedFormats[ext] {
			return nil
		}

		log.Printf("Found book: %s\n", path)

		// Check if book already exists in database
		exists, err := l.bookExists(path)
		if err != nil {
			return err
		}

		if exists {
			// Update existing book if file was modified
			return l.updateBookIfModified(path, info)
		}

		// Add new book
		return l.addBook(path, info)
	})
}

func (l *Library) bookExists(filePath string) (bool, error) {
	var count int
	err := l.DB.QueryRow("SELECT COUNT(*) FROM books WHERE filepath = ?", filePath).Scan(&count)
	return count > 0, err
}

func (l *Library) updateBookIfModified(filePath string, info os.FileInfo) error {
	var dbModTime time.Time
	err := l.DB.QueryRow("SELECT modtime FROM books WHERE filepath = ?", filePath).Scan(&dbModTime)
	if err != nil {
		return err
	}

	if info.ModTime().After(dbModTime) {
		title := l.extractTitle(filePath)
		format := strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), ".")

		_, err = l.DB.Exec(`
			UPDATE books 
			SET title = ?, filesize = ?, modtime = ?, format = ?
			WHERE filepath = ?`,
			title, info.Size(), info.ModTime(), format, filePath)
		return err
	}

	return nil
}

func (l *Library) addBook(filePath string, info os.FileInfo) error {
	title := l.extractTitle(filePath)
	format := strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), ".")

	_, err := l.DB.Exec(`
		INSERT INTO books (filepath, title, filesize, modtime, format)
		VALUES (?, ?, ?, ?, ?)`,
		filePath, title, info.Size(), info.ModTime(), format)
	return err
}

func (l *Library) extractTitle(filePath string) string {
	// Try to get title from metadata using exiftool
	title := zathura.GetFileTitle(filePath)
	
	// If exiftool fails or returns empty, use filename without extension
	if strings.TrimSpace(title) == "" {
		title = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}
	
	return strings.TrimSpace(title)
}

func (l *Library) GetAllBooks() ([]Book, error) {
	rows, err := l.DB.Query(`
		SELECT id, filepath, title, filesize, modtime, format
		FROM books
		ORDER BY title ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.FilePath, &book.Title, &book.FileSize, &book.ModTime, &book.Format)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, rows.Err()
}

func (l *Library) GetBooksByFormat(format string) ([]Book, error) {
	rows, err := l.DB.Query(`
		SELECT id, filepath, title, filesize, modtime, format
		FROM books
		WHERE format = ?
		ORDER BY title ASC`, format)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.FilePath, &book.Title, &book.FileSize, &book.ModTime, &book.Format)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, rows.Err()
}

func (l *Library) SearchBooks(query string) ([]Book, error) {
	searchQuery := "%" + strings.ToLower(query) + "%"
	rows, err := l.DB.Query(`
		SELECT id, filepath, title, filesize, modtime, format
		FROM books
		WHERE LOWER(title) LIKE ? OR LOWER(filepath) LIKE ?
		ORDER BY title ASC`, searchQuery, searchQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.FilePath, &book.Title, &book.FileSize, &book.ModTime, &book.Format)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, rows.Err()
}

func (l *Library) RemoveBook(filePath string) error {
	_, err := l.DB.Exec("DELETE FROM books WHERE filepath = ?", filePath)
	return err
}

func (l *Library) CleanupMissingFiles() error {
	books, err := l.GetAllBooks()
	if err != nil {
		return err
	}

	for _, book := range books {
		if _, err := os.Stat(book.FilePath); os.IsNotExist(err) {
			if err := l.RemoveBook(book.FilePath); err != nil {
				return fmt.Errorf("failed to remove missing book %s: %w", book.FilePath, err)
			}
		}
	}

	return nil
}

func (l *Library) initBookmarkExtractor() error {
	zathuraDbPath, err := zathura.GetDatabasePath()
	if err != nil {
		return fmt.Errorf("failed to get zathura database path: %w", err)
	}

	extractor, err := zathura.NewBookmarkExtractor(zathuraDbPath)
	if err != nil {
		return fmt.Errorf("failed to create bookmark extractor: %w", err)
	}

	l.Extractor = extractor
	return nil
}

func (l *Library) Close() error {
	// Close the bookmark extractor database connection if it exists
	if l.Extractor != nil && l.Extractor.DB != nil {
		l.Extractor.DB.Close()
	}
	return l.DB.Close()
}
