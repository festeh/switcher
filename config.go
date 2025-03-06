package main

// Command struct represents a command that can be executed
type Command struct {
	Name string
}

// GetCommands returns a list of available commands
func GetCommands() []Command {
	return []Command{
		{Name: "firefox"},
		{Name: "vscode"},
	}
}
