package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/factory-cockpit/workspace"
)

type Model struct {
	tasks       []Task
	selectedIdx int
	workspace   *workspace.Workspace
	manager     *workspace.Manager
	spinner     spinner.Model
	quitting    bool
	loading     bool
	loadingMsg  string
	err         error
	msgChan      chan tuiMsg
}

type Task struct {
	ID          string
	Title       string
	Status      string
	IssueID     string
	Branch      string
	WorkspaceID string
	IsActive    bool
}

func NewModel(mgr *workspace.Manager) *Model {
	s := spinner.New()
	s.Tick()
	
	activeTaskId, _ := mgr.LoadState()
	
	return &Model{
		manager: mgr,
		spinner: s,
		tasks: []Task{
			{ID: "1", Title: "Implement workspace manager", Status: "OPEN", IssueID: "5", Branch: "feature/workspace"},
			{ID: "2", Title: "Add auth", Status: "TODO", IssueID: "6", Branch: "feature/auth"},
		},
		msgChan: make(chan tuiMsg, 10),
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "i":
			if m.selectedIdx >= 0 && m.selectedIdx < len(m.tasks) {
				if m.loading {
					return m, nil
				}
				m.loading = true
				m.loadingMsg = fmt.Sprintf("Launching workspace for %s...", m.tasks[m.selectedIdx].Title)
				return m, m.launchWorkspaceCmd(m.tasks[m.selectedIdx].ID, m.selectedIdx)
			}
			return m, nil
		case "x":
			if m.workspace != nil {
				return m, m.closeWorkspaceCmd()
			}
			return m, nil
		case "up":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil
		case "down":
			if m.selectedIdx < len(m.tasks)-1 {
				m.selectedIdx++
			}
			return m, nil
		}
	case spinner.Tick:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tuiMsg:
		m.loading = false
		m.loadingMsg = ""
		m.err = msg.err
		if msg.ws != nil {
			m.workspace = msg.ws
			// Update the active task icon
			for i := range m.tasks {
				if m.tasks[i].ID == msg.ws.TaskID {
					m.tasks[i].IsActive = true
				} else {
					m.tasks[i].IsActive = false
				}
			}
		}
		return m, nil
	}
	return m, nil
}

type tuiMsg struct {
	ws  *workspace.Workspace
	err error
}

func (m *Model) launchWorkspaceCmd(taskId string, idx int) tea.Cmd {
	return func() tea.Msg {
		ws, ok := m.manager.SessionMap[taskId]
		if !ok {
			task := m.tasks[idx]
			ws, err := m.manager.GetOrCreateWorkspace(task.ID, task.IssueID, task.Branch)
			if err != nil {
				// This is a bit of a hack since we are in a Cmd, 
				// but we'll have to handle the error via the msg.
				return tuiMsg{err: err}
			}
		}
		
		m.manager.LaunchWorkspace(ws)
		m.manager.AttachToSession(ws)
		
		_ = m.manager.SaveState(taskId)
		
		return tuiMsg{ws: ws, err: nil}
	}
}

func (m *Model) closeWorkspaceCmd() tea.Cmd {
	return func() tea.Msg {
		if m.workspace != nil {
			m.manager.Detach(m.workspace)
			m.manager.SaveState("")
			m.workspace = nil
			for i := range m.tasks {
				m.tasks[i].IsActive = false
			}
		}
		return tuiMsg{}
	}
}

func (m Model) View() string {
	if m.quitting {
		return "Quitting...\n"
	}

	if m.loading {
		return fmt.Sprintf("\n%s\n\n%s\n\n", m.spinner.View(), m.loadingMsg)
	}

	s := "🏭 DEV FACTORY COCKPIT\n\n"
	if m.workspace != nil {
		s += fmt.Sprintf("Active Workspace: %s\n", m.workspace.TmuxSession)
	} else {
		s += "Active Workspace: None\n"
	}
	s += "\nTasks:\n"
	for i, t := range m.tasks {
		cursor := "  "
		if m.selectedIdx == i {
			cursor = "> "
		}
		icon := ""
		if t.IsActive {
			icon = "▶ "
		}
		s += fmt.Sprintf("%s [%s] %s %s\n", cursor, t.Status, icon, t.Title)
	}
	s += "\n(i) Jump | (x) Close | (q) Quit"
	return s
}
