package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/getlantern/systray"
)

// App struct
type App struct {
	ctx    context.Context
	config Config
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
		fmt.Println(errorMsg)
		os.Exit(1)
	}
	return &App{
		config: config,
	}
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetCommandList returns the list of available commands
func (a *App) GetCommandList() []Command {
	return a.config.Commands
}
