package main

import (
	"context"
	_ "embed"
	"fmt"
	// "os"
	"os/exec"
	"time"

	"github.com/getlantern/systray"
	"github.com/yourusername/switcher/books"
)

// App struct
type App struct {
	ctx       context.Context
	config    Config
	extractor *books.BookmarkExtractor
}

//go:embed assets/letter-s.png
var iconData []byte

func onReady() {
	systray.SetTitle("switcher")
	systray.SetTooltip("switcher")
	systray.AddMenuItem("Quit", "Quit the whole app")
	systray.SetIcon(iconData)
}

func onExit() {
	systray.Quit()
}

// NewApp creates a new App application struct
func NewApp() *App {
	config, err := LoadConfig()
	if err != nil {
		// If config can't be loaded, notify user and exit
		errorMsg := fmt.Sprintf("Failed to load configuration: %v", err)
		exec.Command("notify-send", "Switcher Error", errorMsg).Run()
		// fmt.Println(errorMsg)
		// os.Exit(1)
	}

	app := &App{
		config: config,
	}

	// Initialize the bookmark extractor
	dbPath, err := books.GetDatabasePath()
	if err == nil {
		extractor, err := books.NewBookmarkExtractor(dbPath)
		if err == nil {
			app.extractor = extractor
		} else {
			fmt.Printf("Failed to create bookmark extractor: %v\n", err)
		}
	} else {
		fmt.Printf("Failed to get database path: %v\n", err)
	}

	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go func() {
		// Wait for Wails to fully initialize
		time.Sleep(500 * time.Millisecond)
		systray.Run(onReady, onExit)
	}()
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	// Close the database connection if extractor exists
	if a.extractor != nil && a.extractor.DB != nil {
		a.extractor.DB.Close()
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetCommandList returns the list of available commands
func (a *App) GetCommandList() []Command {
	// Convert map to slice for frontend compatibility
	var commandList []Command
	for _, cmd := range a.config.Commands {
		commandList = append(commandList, cmd)
	}
	return commandList
}

// ExecCommand executes the specified command
func (a *App) ExecCommand(cmd string) error {
	// First, move the switcher window to a workspace silently using hyprctl
	moveCmd := exec.Command("hyprctl", "dispatch", "movetoworkspacesilent", "9,title:switcher.*")
	if err := moveCmd.Run(); err != nil {
		// If hyprctl fails, log the error but continue with the main command
		fmt.Printf("Failed to move window: %v\n", err)
	}

	// Execute the requested command
	command := exec.Command(cmd)
	return command.Start()
}

// GetBookmarks returns the list of bookmarks from Zathura
func (a *App) GetBookmarks() ([]books.BookmarkInfo, error) {
	if a.extractor == nil {
		// If extractor wasn't initialized during startup, try again now
		dbPath, err := books.GetDatabasePath()
		if err != nil {
			return nil, fmt.Errorf("failed to get database path: %w", err)
		}

		extractor, err := books.NewBookmarkExtractor(dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create bookmark extractor: %w", err)
		}
		defer extractor.DB.Close()
		
		return extractor.ExtractBookmarks()
	}
	
	// Use the existing extractor
	bookmarks, err := a.extractor.ExtractBookmarks()
	if err != nil {
		return nil, fmt.Errorf("failed to extract bookmarks: %w", err)
	}

	return bookmarks, nil
}
