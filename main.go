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
	// Skip check if in development mode
	if os.Getenv("WAILS_DEV") != "" {
		fmt.Println("Development mode detected, skipping process check")
		return
	}
	
	currentPID := os.Getpid()
	processes, err := exec.Command("pgrep", "switcher").Output()
	if err != nil {
		return
	}
	runningPIDs := strings.Split(strings.TrimSpace(string(processes)), "\n")
	for _, pidStr := range runningPIDs {
		pid, err := strconv.Atoi(pidStr)
		if err == nil && pid != currentPID {
			
      // HACK: if workspace resolution fails, it selects current one
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

func runDoctorCommand() {
	fmt.Println("Running Switcher doctor...")
	
	// Try to load the configuration
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("❌ Configuration error: %v\n", err)
		os.Exit(1)
	}
	
	// Print configuration details
	fmt.Println("✅ Configuration loaded successfully")
	fmt.Printf("Found %d commands in configuration\n", len(config.Commands))
	
	// List the commands
	if len(config.Commands) > 0 {
		fmt.Println("Commands:")
		for i, cmd := range config.Commands {
			fmt.Printf("  %d. %s\n", i+1, cmd.Name)
		}
	} else {
		fmt.Println("Warning: No commands defined in configuration")
	}
	
	os.Exit(0)
}

func main() {
	// Check for doctor command
	if len(os.Args) > 1 && os.Args[1] == "doctor" {
		runDoctorCommand()
		return
	}
	
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
