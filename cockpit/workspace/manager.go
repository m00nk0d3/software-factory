package workspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	StatusLeftLength = 60
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

type Workspace struct {
	TaskID       string
	IssueID      string
	WorktreePath string
	Branch       string
	TmuxSession  string
	Active       bool
}

type Manager struct {
	SandcastleBaseURL string
	SessionMap        map[string]*Workspace
	StateFile         string
	mu                sync.RWMutex
}

func NewManager(baseURL string) *Manager {
	return &Manager{
		SandcastleBaseURL: baseURL,
		SessionMap:        make(map[string]*Workspace),
		StateFile:         "workspace_state.json",
	}
}

// SaveState saves the currently active workspace TaskID to a file
func (m *Manager) SaveState(taskId string) error {
	data, _ := json.Marshal(map[string]string{"active_task_id": taskId})
	tmpFile := m.StateFile + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmpFile, m.StateFile)
}

// LoadState loads the active task ID from the file
func (m *Manager) LoadState() (string, error) {
	data, err := os.ReadFile(m.StateFile)
	if err != nil {
		return "", err
	}
	var state map[string]string
	if err := json.Unmarshal(data, &state); err != nil {
		return "", err
	}
	return state["active_task_id"], nil
}

// CreateWorkspace calls the Sandcastle Bridge to create a new worktree
func (m *Manager) CreateWorkspace(taskId, issueId, branch string) (*Workspace, error) {
	payload := map[string]string{
		"taskId":  taskId,
		"issueId": issueId,
		"branch":  branch,
	}
	
	body, _ := json.Marshal(payload)
	resp, err := httpClient.Post(fmt.Sprintf("%s/api/worktree/create", m.SandcastleBaseURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create worktree: %d", resp.StatusCode)
	}

	var result struct {
		WorktreePath string `json:"worktreePath"`
		Branch       string `json:"branch"`
		TaskId       string `json:"taskId"`
		Status       string `json:"status"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	ws := &Workspace{
		TaskID:       result.TaskId,
		IssueID:      issueId,
		WorktreePath: result.WorktreePath,
		Branch:       result.Branch,
		TmuxSession:  fmt.Sprintf("factory-%s", issueId),
		Active:       true,
	}

	m.mu.Lock()
	m.SessionMap[taskId] = ws
	m.mu.Unlock()
	
	if err := m.SaveState(taskId); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save state: %v\n", err)
	}
	return ws, nil
}

// LaunchWorkspace initializes the tmux session for the given workspace
func (m *Manager) LaunchWorkspace(ws *Workspace) error {
	// Create new session in detached mode
	_, err := exec.Command("tmux", "new-session", "-d", "-s", ws.TmuxSession, "-c", ws.WorktreePath).Run()
	if err != nil {
		return fmt.Errorf("failed to create tmux session: %w", err)
	}

	// Split window vertically (50/50)
	_, err = exec.Command("tmux", "split-window", "-h", "-t", ws.TmuxSession+":0", "-c", ws.WorktreePath).Run()
	if err != nil {
		return fmt.Errorf("failed to split tmux window: %w", err)
	}

	// Left pane: Neovim
	_, err = exec.Command("tmux", "send-keys", "-t", ws.TmuxSession+":0.0", "nvim .", "Enter").Run()
	if err != nil {
		return fmt.Errorf("failed to send keys to left pane: %w", err)
	}

	// Right pane: Pi Agent
	piCmd := fmt.Sprintf("pi --worktree-mode --task-id=%s", ws.TaskID)
	_, err = exec.Command("tmux", "send-keys", "-t", ws.TmuxSession+":0.1", piCmd, "Enter").Run()
	if err != nil {
		return fmt.Errorf("failed to send keys to right pane: %w", err)
	}

	// Set pane titles
	if _, err := exec.Command("tmux", "select-pane", "-t", ws.TmuxSession+":0.0", "-T", "NEOVIM").Run(); err != nil {
		return fmt.Errorf("failed to set NEOVIM title: %w", err)
	}
	if _, err := exec.Command("tmux", "select-pane", "-t", ws.TmuxSession+":0.1", "-T", "PI AGENT").Run(); err != nil {
		return fmt.Errorf("failed to set PI AGENT title: %w", err)
	}

	// Configure status bar
	statusLeft := fmt.Sprintf("🏭 Factory | 📂 %s | 🌿 %s", ws.IssueID, ws.Branch)
	if _, err := exec.Command("tmux", "set-option", "-t", ws.TmuxSession, "status-left", statusLeft).Run(); err != nil {
		return fmt.Errorf("failed to set status-left: %w", err)
	}
	if _, err := exec.Command("tmux", "set-option", "-t", ws.TmuxSession, "status-left-length", fmt.Sprintf("%d", StatusLeftLength)).Run(); err != nil {
		return fmt.Errorf("failed to set status-left-length: %w", err)
	}
	if _, err := exec.Command("tmux", "set-option", "-t", ws.TmuxSession, "status-right", "[Ctrl+B D] Exit to Cockpit").Run(); err != nil {
		return fmt.Errorf("failed to set status-right: %w", err)
	}

	// Focus on Neovim pane
	if _, err := exec.Command("tmux", "select-pane", "-t", ws.TmuxSession+":0.0").Run(); err != nil {
		return fmt.Errorf("failed to focus Neovim pane: %w", err)
	}

	return nil
}

// AttachToSession blocks until the user detaches from the tmux session
func (m *Manager) AttachToSession(ws *Workspace) error {
	cmd := exec.Command("tmux", "attach-session", "-t", ws.TmuxSession)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Detach returns from the attached session
func (m *Manager) Detach(ws *Workspace) error {
	// Note: tmux attach-session blocks the process. Detachment is typically 
	// handled by the user via keyboard shortcuts (e.g., Ctrl+B D).
	// This method is kept for interface completeness.
	return nil
}

// ListWorktrees fetches all active worktrees from the Sandcastle Bridge
func (m *Manager) ListWorktrees() ([]Workspace, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	resp, err := httpClient.Get(fmt.Sprintf("%s/api/worktree/list", m.SandcastleBaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list worktrees: %d", resp.StatusCode)
	}

	var results struct {
		Worktrees []Workspace `json:"worktrees"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results.Worktrees, nil
}