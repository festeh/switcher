package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	// Check if another instance of switcher is already running
	currentPID := os.Getpid()
	processes, err := exec.Command("pgrep", "switcher").Output()
	if err == nil {
		runningPIDs := strings.Split(strings.TrimSpace(string(processes)), "\n")
		for _, pidStr := range runningPIDs {
			pid, err := strconv.Atoi(pidStr)
			if err == nil && pid != currentPID {
				// Another instance is running, send notification and exit
				exec.Command("notify-send", "Switcher", "Switcher is already running").Run()
				fmt.Println("Switcher is already running. Exiting.")
				os.Exit(0)
			}
		}
	}
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "switcher",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 245, G: 247, B: 250, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
