package library

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"switcher/foliate"
	"switcher/util"
	"switcher/zathura"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Book struct {
	FilePath string `json:"filepath"`
	Title    string `json:"title"`
	Format   string `json:"format"`
	Page     int    `json:"page,omitempty"`
}

type Library struct {
	DB      *sql.DB
	Zathura *zathura.Zathura
	Foliate *foliate.Foliate
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

	zat, err := zathura.NewZathura()
	if err != nil {
		return nil, err
	}

	foli, err := foliate.NewFoliate()
	if err != nil {
		return nil, err
	}

	library.Zathura = zat
	library.Foliate = foli

	return library, nil
}

func (l *Library) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS books (
		filepath TEXT UNIQUE NOT NULL,
		title TEXT NOT NULL,
		format TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_books_filepath ON books(filepath);
	`
	_, err := l.DB.Exec(query)
	return err
}

func (l *Library) ResetDatabase() error {
	query := `DROP TABLE IF EXISTS books;`
	_, err := l.DB.Exec(query)
	if err != nil {
		return err
	}
	return l.initSchema()
}

func getIgnoredDirs(rootDir string) (map[string]struct{}, error) {
	ignoredDirs := make(map[string]struct{})
	ignoreFilePath := filepath.Join(rootDir, ".ignore")

	file, err := os.Open(ignoreFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ignoredDirs, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			ignoredDirs[line] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ignoredDirs, nil
}

func (l *Library) ScanDirectory(rootDir string) error {
	supportedFormats := map[string]bool{
		".pdf":  true,
		".epub": true,
		".fb2":  true,
	}

	ignoredDirs, err := getIgnoredDirs(rootDir)
	if err != nil {
		return fmt.Errorf("error reading .ignore file: %w", err)
	}

	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if filepath.Dir(path) == rootDir {
				dirName := info.Name()
				if _, ok := ignoredDirs[dirName]; ok {
					log.Printf("Ignoring directory: %s\n", path)
					return filepath.SkipDir
				}
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !supportedFormats[ext] {
			return nil
		}

		log.Printf("Found book: %s\n", path)

		exists, err := l.bookExists(path)
		if err != nil {
			return err
		}

		if !exists {
			return l.addBook(path)
		}
		return nil
	})
}

func (l *Library) bookExists(filePath string) (bool, error) {
	var count int
	err := l.DB.QueryRow("SELECT COUNT(*) FROM books WHERE filepath = ?", filePath).Scan(&count)
	return count > 0, err
}

func (l *Library) addBook(filePath string) error {
	title := l.extractTitle(filePath)
	format := strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), ".")

	_, err := l.DB.Exec(`
		INSERT INTO books (filepath, title, format)
		VALUES (?, ?, ?)`,
		filePath, title, format)
	return err
}

func (l *Library) extractTitle(filePath string) string {
	cmd := exec.Command("exiftool", "-s", "-s", "-s", "-Title", filePath)
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing exiftool on %s: %v", filePath, err)
		return filepath.Base(filePath) // Return filename as fallback
	}
	title := string(out)

	if strings.TrimSpace(title) == "" {
		cmd := exec.Command("exiftool", filePath)
		fullOut, err := cmd.Output()
		if err != nil {
			log.Printf("Error executing exiftool for full output on %s: %v", filePath, err)
			// Fallback to filename without extension
			title = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
		} else {
			scanner := bufio.NewScanner(strings.NewReader(string(fullOut)))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, "Book-title:") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) > 1 {
						title = parts[1]
						break // Found it, no need to scan further
					}
				}
			}
		}
	}

	// If exiftool fails or returns empty, use filename without extension
	if strings.TrimSpace(title) == "" {
		title = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}

	return strings.TrimSpace(title)
}

func (l *Library) SearchBooks(term string) ([]Book, error) {
	allBooks, err := l.GetAllBooks()
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(term) == "" {
		return allBooks, nil
	}

	var titles []string
	for _, book := range allBooks {
		titles = append(titles, book.Title)
	}

	ranks := fuzzy.RankFindFold(term, titles)
	sort.Sort(ranks)

	var foundBooks []Book
	for _, rank := range ranks {
		foundBooks = append(foundBooks, allBooks[rank.OriginalIndex])
	}

	return foundBooks, nil
}

func (l *Library) GetAllBooks() ([]Book, error) {
	zathuraMap, err := l.Zathura.GetAllKnownBooks()
	if err != nil {
		log.Printf("could not get zathura books, continuing without them: %v", err)
	}

	foliateMap, err := l.Foliate.GetAllKnownBooks()
	if err != nil {
		log.Printf("could not get foliate books, continuing without them: %v", err)
	}

	rows, err := l.DB.Query(`
		SELECT filepath, title, format
		FROM books
		ORDER BY title ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.FilePath, &book.Title, &book.Format)
		if err != nil {
			return nil, err
		}

		if zathuraBook, ok := zathuraMap[book.FilePath]; ok {
			book.Page = zathuraBook.Page
		}

		if foliateBook, ok := foliateMap[book.FilePath]; ok {
			book.Title = foliateBook.Title
			book.Page = foliateBook.Page
		}

		books = append(books, book)
	}

	return books, rows.Err()
}

func (l *Library) GetBooksByFormat(format string) ([]Book, error) {
	rows, err := l.DB.Query(`
		SELECT filepath, title,  format
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
		err := rows.Scan(&book.FilePath, &book.Title, &book.Format)
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

func (l *Library) Close() error {
	// Close the bookmark extractor database connection if it exists
	if l.Zathura != nil && l.Zathura.DB != nil {
		l.Zathura.DB.Close()
	}
	return l.DB.Close()
}
