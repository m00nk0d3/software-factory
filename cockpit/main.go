package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/yourusername/factory-cockpit/workspace"
)

func main() {
	// Configuration
	sandcastleBaseURL := os.Getenv("SANDCASTLE_URL")
	if sandcastleBaseURL == "" {
		sandcastleBaseURL = "http://localhost:3001"
	}

	stateFilePath := os.Getenv("STATE_FILE_PATH")
	if stateFilePath == "" {
		configDir, err := os.UserConfigDir()
		if err == nil {
			stateFilePath = fmt.Sprintf("%s/software-factory/workspace_state.json", configDir)
		} else {
			stateFilePath = "workspace_state.json"
		}
	}

	// Initialize Bridge and Manager
	bridge := workspace.NewHttpWorkspaceBridge(sandcastleBaseURL, nil)
	mgr := workspace.NewManager(bridge, stateFilePath)
	
	// Load persistence state
	activeTaskId, err := mgr.LoadState()
	if err == nil && activeTaskId != "" {
		fmt.Printf("Restoring session for task: %s\n", activeTaskId)
		// We'll use this to flag the correct task in the UI
	}

	p := tui.NewModel(mgr)
	
	// Inject the active task ID into the model
	if activeTaskId != "" {
		for i := range p.tasks {
			if p.tasks[i].ID == activeTaskId {
				p.tasks[i].IsActive = true
			}
		}
	}

	if _, err := bubbletea.Run(p); err != nil {
		fmt.Printf("Error running app: %v\n", err)
		os.Exit(1)
	}
}
