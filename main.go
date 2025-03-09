package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

func checkAlreadyRuns() {
	currentPID := os.Getpid()
	processes, err := exec.Command("pgrep", "switcher").Output()
	if err != nil {
		return
	}
	runningPIDs := strings.Split(strings.TrimSpace(string(processes)), "\n")
	for _, pidStr := range runningPIDs {
		pid, err := strconv.Atoi(pidStr)
		if err == nil && pid != currentPID {
			// Another instance is running, try to move its window to current workspace
			// First get the current workspace
			workspaceCmd := exec.Command("hyprctl", "activeworkspace", "-j")
			workspaceOutput, err := workspaceCmd.Output()
			if err != nil {
				fmt.Println("Failed to get current workspace:", err)
				os.Exit(0)
			}
			
			// Move the window to the current workspace
			moveCmd := exec.Command("hyprctl", "dispatch", "movetoworkspace", "special:current,title:switcher")
			err = moveCmd.Run()
			if err != nil {
				fmt.Println("Failed to move switcher window:", err)
				os.Exit(0)
			}
			
			// Then focus the window
			focusCmd := exec.Command("hyprctl", "dispatch", "focuswindow", "title:switcher")
			err = focusCmd.Run()
			if err != nil {
				fmt.Println("Failed to focus switcher window:", err)
				os.Exit(0)
			}
			
			fmt.Println("Moved and focused existing switcher window. Exiting.")
			os.Exit(0)
		}
	}
}

func main() {
	checkAlreadyRuns()
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
