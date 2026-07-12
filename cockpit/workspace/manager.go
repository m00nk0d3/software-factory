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

// WorkspaceBridge defines the interface for interacting with the Sandcastle Bridge.
type WorkspaceBridge interface {
	CreateWorktree(taskId, issueId, branch string) (*Workspace, error)
	ListWorktrees() ([]Workspace, error)
}

// HttpWorkspaceBridge is the production implementation of WorkspaceBridge.
type HttpWorkspaceBridge struct {
	BaseURL string
	Client  *http.Client
}

func NewHttpWorkspaceBridge(baseURL string, client *http.Client) *HttpWorkspaceBridge {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &HttpWorkspaceBridge{
		BaseURL: baseURL,
		Client:  client,
	}
}

func (b *HttpWorkspaceBridge) CreateWorktree(taskId, issueId, branch string) (*Workspace, error) {
	payload := map[string]string{
		"taskId":  taskId,
		"issueId": issueId,
		"branch":  branch,
	}

	body, _ := json.Marshal(payload)
	resp, err := b.Client.Post(fmt.Sprintf("%s/api/worktree/create", b.BaseURL), "application/json", bytes.NewBuffer(body))
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

	return &Workspace{
		TaskID:       result.TaskId,
		IssueID:      issueId,
		WorktreePath: result.WorktreePath,
		Branch:       result.Branch,
		TmuxSession:  fmt.Sprintf("factory-%s", issueId),
		Active:       true,
	}, nil
}

func (b *HttpWorkspaceBridge) ListWorktrees() ([]Workspace, error) {
	resp, err := b.Client.Get(fmt.Sprintf("%s/api/worktree/list", b.BaseURL))
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

type Workspace struct {
	TaskID       string
	IssueID      string
	WorktreePath string
	Branch       string
	TmuxSession  string
	Active       bool
}

type Manager struct {
	Bridge       WorkspaceBridge
	SessionMap   map[string]*Workspace
	StateFile    string
	mu           sync.RWMutex
}

func NewManager(bridge WorkspaceBridge, stateFilePath string) *Manager {
	return &Manager{
		Bridge:     bridge,
		SessionMap: make(map[string]*Workspace),
		StateFile:  stateFilePath,
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

// GetOrCreateWorkspace returns an existing workspace or creates a new one
func (m *Manager) GetOrCreateWorkspace(taskId, issueId, branch string) (*Workspace, error) {
	m.mu.RLock()
	ws, ok := m.SessionMap[taskId]
	m.mu.RUnlock()

	if ok {
		return ws, nil
	}

	ws, err := m.Bridge.CreateWorktree(taskId, issueId, branch)
	if err != nil {
		return nil, err
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
	// 1. Create new session in detached mode
	// Using -c to set the initial directory correctly
	_, err := exec.Command("tmux", "new-session", "-d", "-s", ws.TmuxSession, "-c", ws.WorktreePath).Run()
	if err != nil {
		return fmt.Errorf("failed to create tmux session: %w", err)
	}

	// 2. Split window vertically
	// We target the session name directly to ensure we are operating on the correct window
	_, err = exec.Command("tmux", "split-window", "-h", "-t", ws.TmuxSession, "-c", ws.WorktreePath).Run()
	if err != nil {
		return fmt.Errorf("failed to split tmux window: %w", err)
	}

	// 3. Target panes explicitly to avoid ambiguity
	// Pane 0.0 is left, 0.1 is right (standard split behavior)
	leftPane := ws.TmuxSession + ":0.0"
	rightPane := ws.TmuxSession + ":0.1"

	// Left pane: Neovim
	_, err = exec.Command("tmux", "send-keys", "-t", leftPane, "nvim .", "Enter").Run()
	if err != nil {
		return fmt.Errorf("failed to send keys to left pane: %w", err)
	}

	// Right pane: Pi Agent
	piCmd := fmt.Sprintf("pi --worktree-mode --task-id=%s", ws.TaskID)
	_, err = exec.Command("tmux", "send-keys", "-t", rightPane, piCmd, "Enter").Run()
	if err != nil {
		return fmt.Errorf("failed to send keys to right pane: %w", err)
	}

	// 4. Set pane titles and status
	titles := map[string]string{
		leftPane:  "NEOVIM",
		rightPane: "PI AGENT",
	}
	for target, title := range titles {
		if _, err := exec.Command("tmux", "select-pane", "-t", target, "-T", title).Run(); err != nil {
			return fmt.Errorf("failed to set title for %s: %w", target, err)
		}
	}

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
	if _, err := exec.Command("tmux", "select-pane", "-t", leftPane).Run(); err != nil {
		return fmt.Errorf("failed to focus Neovim pane: %w", err)
	}

	return nil
}

// ListWorktrees fetches all active worktrees from the bridge
func (m *Manager) ListWorktrees() ([]Workspace, error) {
	return m.Bridge.ListWorktrees()
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
	return nil
}
