# Dev Factory Implementation Guide - Complete Architecture

**Date:** 2026-07-12  
**Status:** Design Complete, Ready for Implementation  
**Hardware:** 32GB RAM, RTX 5070 Ti (16GB VRAM)

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Original Vision](#original-vision)
3. [Architectural Evolution](#architectural-evolution)
4. [Final Tech Stack](#final-tech-stack)
5. [The Factory Cockpit](#the-factory-cockpit)
6. [Integrated Workspace](#integrated-workspace)
7. [Core Architecture](#core-architecture)
8. [Implementation Details](#implementation-details)
9. [Hardware Optimization](#hardware-optimization)
10. [Implementation Roadmap](#implementation-roadmap)
11. [Code Reference](#code-reference)
12. [Next Steps](#next-steps)

---

## Executive Summary

### What We're Building

A **fully autonomous, local-first software development factory** with a beautiful terminal-based cockpit interface. The system handles the entire SDLC from requirements gathering to PR creation, with an integrated development environment that combines Neovim and Pi Agent in split-screen workspaces.

### Key Principles

- **100% Free & Local**: No API costs, no vendor lock-in
- **Agent Specialization**: Each agent type has a specific role (PM, Architect, Coder, Reviewer)
- **Event-Driven**: Async message passing for scalability and resilience
- **Worktree Isolation**: Every task gets its own git worktree (managed by Sandcastle)
- **TDD-First**: Tests are generated before implementation
- **Terminal-First UI**: Beautiful Go TUI cockpit with integrated workspaces
- **Seamless IDE Experience**: Neovim + Pi Agent side-by-side in tmux
- **Observable**: Full tracing, metrics, and logging

### System Components

```
┌─────────────────────────────────────────────────────────┐
│            Factory Cockpit (Go TUI)                     │
│  Beautiful terminal dashboard for monitoring & control  │
│  Press 'i' on any task → Jump into workspace           │
└────────────┬────────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────────────────┐
│         Integrated Workspace (Tmux Split)               │
│  ┌──────────────────┬──────────────────────────────┐   │
│  │   NEOVIM         │   PI AGENT                   │   │
│  │   Code Editor    │   AI Assistant               │   │
│  └──────────────────┴──────────────────────────────┘   │
└────────────┬────────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────────────────┐
│         Python Orchestration Layer                      │
│  - Event-driven agent coordination                      │
│  - State management (PostgreSQL + Redis)                │
│  - Workflow execution                                   │
└────────────┬────────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────────────────┐
│         Sandcastle (Worktree Management)                │
│  - Git worktree isolation per task                      │
│  - Branch strategies                                    │
│  - Lifecycle management                                 │
└─────────────────────────────────────────────────────────┘
```

### Expected Performance

With your hardware (32GB RAM, RTX 5070 Ti):
- **Throughput**: 10-15 tasks/hour (sequential), 20-30 tasks/hour (parallel)
- **Latency**: Simple tasks in 10-30 seconds, complex features in 5-10 minutes
- **Cost**: ~$35/month electricity vs $500-2000/month for cloud APIs

---

## Original Vision

Based on `dev-workflow.md`, the goal was to transition from "AI-assisted coding" to "Agentic Development" with these workstreams:

### A. Product & Strategy (The "Think" Tank)
- **PM Agent**: Extract requirements, identify edge cases
- **Architect Agent**: Generate PRDs, TRDs, system designs
- **Analyst Agent**: Decompose PRDs into granular tasks

### B. Implementation (The "Build" Squad)
- **Router Agent**: Route tasks to specialist agents
- **Specialist Agents**: Framework-specific (React, Python, Go, etc.)
- **TDD Specialist**: Generate tests before code
- **Sandbox Executor**: Run code in isolated containers

### C. Quality & Governance (The "Gatekeepers")
- **Reviewer Agent**: Deep code review
- **Security Agent**: Scan for vulnerabilities
- **QA Specialist**: Validate Definition of Done

### D. Logistics & DevOps (The "Nervous System")
- **Git Specialist**: All git operations
- **PR Manager**: Create and manage PRs
- **Workspace Manager**: Manage tmux/neovim sessions
- **State Registry**: Track all active worktrees

---

## Architectural Evolution

### Phase 1: Cloud-First Design (Initial)

Original plan used:
- **Temporal** for orchestration
- **NATS** for event bus
- **Anthropic/OpenAI APIs** for LLM inference
- **E2B** for sandboxing
- **Kubernetes** for deployment

**Problem**: Monthly costs of $500-2000+, vendor lock-in, external dependencies.

### Phase 2: Local-First Revolution

Rebuilt with 100% free, self-hosted tools:
- **Python asyncio** replaces Temporal
- **Redis Streams** replaces NATS
- **Ollama** for local LLM inference → **Changed to LM Studio**
- **Docker** for sandboxing
- **Docker Compose** replaces Kubernetes

### Phase 3: Sandcastle Integration

Discovered **Sandcastle** - a TypeScript library specifically designed for orchestrating AI agents in isolated git worktrees:
- Automatic worktree creation/cleanup
- Branch strategies (head, merge-to-head, branch)
- Docker/Podman sandbox providers
- Lifecycle hooks

### Phase 4: Cockpit UI (Final)

Added a **Go TUI cockpit** with integrated workspaces:
- **Bubbletea** framework for beautiful terminal UI
- **Tmux** integration for split Neovim + Pi workspaces
- **Real-time monitoring** of agents, tasks, and resources
- **Keyboard-driven** workflow

---

## Final Tech Stack

### Human Interface (The Cockpit)
```yaml
cockpit:
  language: Go 1.22+
  framework: Bubbletea (Charm.sh)
  styling: Lipgloss
  components: Bubbles
  terminal: Any modern terminal with 256 color support

workspace:
  multiplexer: Tmux
  editor: Neovim
  ai_assistant: Pi Agent
  layout: Side-by-side split (50/50 or custom)
```

### Worktree Management
```yaml
worktree_engine: Sandcastle (TypeScript)
integration: HTTP API bridge to Python orchestrator
branch_strategy: "branch" (explicit branch per task)
isolation: Docker per worktree
```

### LLM Inference
```yaml
engine: LM Studio
api: OpenAI-compatible (http://localhost:1234/v1)
models:
  primary: deepseek-coder-v2-16b-q4_K_M (9GB VRAM)
  fast: qwen2.5-coder-7b-q4_K_M (4.5GB VRAM)
  reasoning: qwen2.5-14b-q4_K_M (8GB VRAM)
  review: qwen2.5-coder-14b-q4_K_M (8GB VRAM)
  embeddings: nomic-embed-text-v1.5 (250MB VRAM)
```

### Infrastructure
```yaml
orchestration: Python asyncio + custom workflow engine
event_bus: Redis Streams
database: PostgreSQL 15
cache: Redis 7
vector_db: ChromaDB (embedded)
sandbox: Docker (via Sandcastle)
```

### Observability
```yaml
tracing: Jaeger (self-hosted)
metrics: Prometheus + Grafana
logs: Structured logs to file
apm: OpenTelemetry SDK
```

---

## The Factory Cockpit

### Overview

A beautiful, keyboard-driven TUI built with Go and Bubbletea that serves as your mission control for the entire dev factory.

### Main Dashboard View

```
╔══════════════════════════════════════════════════════════════════════════════╗
║  🏭 DEV FACTORY COCKPIT                    [RTX 5070 Ti: 9.2/16GB] [CPU: 45%]║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  📊 OVERVIEW                                                                 ║
║  ┌────────────────────────────────────────────────────────────────────────┐ ║
║  │  Active Tasks: 3          Completed Today: 12      Failed: 1           │ ║
║  │  Active Agents: 5         Queue Depth: 7           Avg Time: 4.2min    │ ║
║  └────────────────────────────────────────────────────────────────────────┘ ║
║                                                                              ║
║  🤖 ACTIVE AGENTS                                                            ║
║  ┌────────────────────────────────────────────────────────────────────────┐ ║
║  │ ⚡ Specialist #1  │ issue-123  │ Implementing...  │ ████████░░ 80%     │ ║
║  │ 🧪 TDD Agent      │ issue-456  │ Writing tests    │ ███░░░░░░░ 30%     │ ║
║  │ 👁️  Reviewer      │ issue-789  │ Code review      │ █████████░ 90%     │ ║
║  └────────────────────────────────────────────────────────────────────────┘ ║
║                                                                              ║
║  📋 TASK QUEUE                                                               ║
║  ┌────────────────────────────────────────────────────────────────────────┐ ║
║  │ → [▶] #123 Add dark mode toggle              [i] Jump In               │ ║
║  │   [ ] #456 Fix login bug                     [PENDING]                 │ ║
║  │   [ ] #789 Update README                     [PENDING]                 │ ║
║  └────────────────────────────────────────────────────────────────────────┘ ║
║                                                                              ║
║  🔔 APPROVAL QUEUE (2)                                                       ║
║  ┌────────────────────────────────────────────────────────────────────────┐ ║
║  │ ⚠️  PR #123: Add dark mode          [A]pprove [R]eject [V]iew          │ ║
║  │ ⚠️  PR #456: Fix login bug          [A]pprove [R]eject [V]iew          │ ║
║  └────────────────────────────────────────────────────────────────────────┘ ║
║                                                                              ║
║  💬 LIVE LOG STREAM                                                          ║
║  ┌────────────────────────────────────────────────────────────────────────┐ ║
║  │ [14:32:45] Specialist #1: Generated 3 files for issue-123              │ ║
║  │ [14:32:46] Sandbox: Running tests in container-issue-123...            │ ║
║  │ [14:32:48] ✓ Tests passed (12/12)                                      │ ║
║  │ [14:32:49] Reviewer: Starting code review...                           │ ║
║  └────────────────────────────────────────────────────────────────────────┘ ║
║                                                                              ║
║  [T]asks [A]gents [W]orktrees [L]ogs [M]etrics [S]ettings [Q]uit          ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Features

**1. Real-Time Monitoring**
- Live agent status with progress bars
- Task queue with status indicators
- Resource usage (VRAM, CPU, RAM)
- Log streaming

**2. Task Management**
- Browse all tasks
- Filter by status/priority
- Jump into any task's workspace with `i`
- Close active workspaces with `x`

**3. Approval Queue**
- Quick approve/reject PRs
- View diffs
- Add comments
- Batch operations

**4. Agent Inspector**
- View agent activity
- Check token usage
- See conversation history
- Kill/restart agents

**5. Worktree Browser**
- List all active worktrees
- Show workspace status
- Quick cleanup
- Detach/reattach sessions

**6. Metrics Dashboard**
- Token usage over time
- Task completion rate
- Agent performance stats
- System resource graphs

### Technology

**Built with:**
- **Bubbletea**: Elm-architecture TUI framework
- **Lipgloss**: Styling and layout
- **Bubbles**: Pre-built components (tables, spinners, progress bars)
- **WebSocket client**: Real-time updates from Python backend

**Why Go + Bubbletea?**
- ✅ Single binary (easy distribution)
- ✅ Fast and responsive
- ✅ Beautiful by default
- ✅ Keyboard-driven workflow
- ✅ SSH-friendly (no X11 needed)
- ✅ Used by K9s, lazygit, gh-dash

---

## Integrated Workspace

### The Core Experience

When you press `i` on a task in the cockpit, you jump into a **fully integrated development environment**:

```
┌────────────────────────────────────────────────────────────────────────┐
│  WORKSPACE: issue-123 (feature/dark-mode)        [Ctrl+B D] Exit      │
├─────────────────────────────────┬──────────────────────────────────────┤
│  NEOVIM                         │  PI AGENT                            │
│  ~/worktrees/issue-123/         │  Connected to: issue-123             │
│  ┌─────────────────────────┐   │  ┌──────────────────────────────┐    │
│  │ src/settings.tsx        │   │  │ 🤖 Task Context Loaded       │    │
│  │                         │   │  │                              │    │
│  │ export function Settings│   │  │ Task: Add dark mode toggle   │    │
│  │   const [dark, setDark] │   │  │ Branch: feature/dark-mode    │    │
│  │   ...                   │   │  │                              │    │
│  │                         │   │  │ How can I help?              │    │
│  │                         │   │  │                              │    │
│  │ [INSERT] -- 12/45 --    │   │  │ > _                          │    │
│  └─────────────────────────┘   │  └──────────────────────────────┘    │
│                                 │                                      │
├─────────────────────────────────┴──────────────────────────────────────┤
│  [Ctrl+B ←/→] Switch Panes  [Ctrl+B D] Detach  [Ctrl+B Z] Zoom       │
└────────────────────────────────────────────────────────────────────────┘
```

### How It Works

**1. Worktree Creation (Sandcastle)**
```typescript
// Sandcastle creates isolated worktree
await run({
  agent: customAgent,
  sandbox: docker(),
  branchStrategy: {
    type: "branch",
    branch: `feature/issue-${issueId}`
  },
  // Worktree is created at ~/worktrees/issue-123/
});
```

**2. Tmux Session Launch (Go Cockpit)**
```go
// Cockpit creates tmux session with split layout
func (m *Manager) LaunchWorkspace(ws *Workspace) error {
    sessionName := fmt.Sprintf("factory-%s", ws.IssueID)
    
    // Create session
    exec.Command("tmux", "new-session", "-d", "-s", sessionName,
        "-c", ws.WorktreePath).Run()
    
    // Split window vertically (50/50)
    exec.Command("tmux", "split-window", "-h", 
        "-t", sessionName).Run()
    
    // Left pane: Neovim
    exec.Command("tmux", "send-keys", "-t", sessionName+":0.0",
        "nvim .", "Enter").Run()
    
    // Right pane: Pi Agent
    exec.Command("tmux", "send-keys", "-t", sessionName+":0.1",
        fmt.Sprintf("pi --worktree-mode --task-id=%s", ws.TaskID),
        "Enter").Run()
    
    // Attach (blocks until user detaches)
    cmd := exec.Command("tmux", "attach-session", "-t", sessionName)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    return cmd.Run()
}
```

**3. Pi Agent Worktree Mode**
```bash
# Pi Agent launched with task context
pi --worktree-mode \
    --task-id=123 \
    --worktree=/home/user/worktrees/issue-123

# Pi automatically:
# - Loads task context from factory API
# - Sets CWD to worktree
# - Shows task-specific prompt
# - Tracks changes and reports back
```

### Layout Options

**Side-by-Side (Default)**
```
┌─────────────────┬─────────────────┐
│                 │                 │
│    NEOVIM       │   PI AGENT      │
│     (50%)       │     (50%)       │
│                 │                 │
└─────────────────┴─────────────────┘
```

**Three-Pane**
```
┌─────────────────┬─────────────────┐
│                 │                 │
│    NEOVIM       │   PI AGENT      │
│     (60%)       │     (40%)       │
│                 │                 │
├─────────────────┴─────────────────┤
│        TERMINAL (30%)             │
│   (for tests, git, etc.)          │
└───────────────────────────────────┘
```

**Custom Ratios**
```yaml
# config/workspace-layouts.yaml
layouts:
  default:
    type: vertical_split
    left_width: 50%
    panes:
      - neovim
      - pi_agent
  
  code_heavy:
    type: vertical_split
    left_width: 70%
    panes:
      - neovim
      - pi_agent
  
  three_pane:
    type: custom
    panes:
      - name: neovim
        position: top_left
        width: 60%
        height: 70%
      - name: pi_agent
        position: top_right
        width: 40%
        height: 70%
      - name: terminal
        position: bottom
        width: 100%
        height: 30%
```

### Workspace Features

**1. Persistent Sessions**
- Workspaces persist after detach
- Press `i` again to reattach
- State is preserved (open files, cursor position, Pi history)

**2. Context Awareness**
- Pi Agent knows the task context
- Neovim opens relevant files
- Tests are pre-loaded
- Git state is ready

**3. Status Bar**
```
🏭 Factory | 📂 issue-123 | 🌿 feature/dark-mode | Ctrl+B D to exit
```

**4. Quick Actions**
```
Ctrl+B ←/→  - Switch between panes
Ctrl+B D    - Detach (return to cockpit)
Ctrl+B Z    - Zoom current pane (fullscreen)
Ctrl+B C    - Create new window
Ctrl+B [    - Scroll mode (for logs)
```

### User Workflow

```bash
# 1. Start cockpit
$ factory-cockpit

# 2. Browse tasks, navigate with ↑/↓
# 3. Press 'i' on a task
#    → Tmux session launches
#    → Neovim opens on left
#    → Pi Agent opens on right

# 4. Work on the task
#    - Edit code in Neovim
#    - Ask Pi for help
#    - Run tests
#    - Commit changes

# 5. Press Ctrl+B D to detach
#    → Returns to cockpit
#    → Workspace still running (▶ icon shown)

# 6. Later, press 'i' again on same task
#    → Reattaches to existing workspace
#    → Everything preserved
```

---

## Core Architecture

### System Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                      Your Local Machine                             │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐ │
│  │           Factory Cockpit (Go + Bubbletea)                    │ │
│  │  - Dashboard, task browser, approval queue                    │ │
│  │  - Launches integrated workspaces (tmux + nvim + pi)          │ │
│  │  - Real-time monitoring via WebSocket                         │ │
│  └──────────────────────────┬────────────────────────────────────┘ │
│                             │                                       │
│  ┌──────────────────────────┴────────────────────────────────────┐ │
│  │          Sandcastle Service (Node.js/TypeScript)              │ │
│  │  - Worktree creation/management                               │ │
│  │  - Branch strategies                                          │ │
│  │  - Docker sandbox provider                                    │ │
│  └──────────────────────────┬────────────────────────────────────┘ │
│                             │                                       │
│  ┌──────────────────────────┴────────────────────────────────────┐ │
│  │       Python Orchestrator (Event-Driven Core)                 │ │
│  │  - Event bus (Redis Streams)                                  │ │
│  │  - State management (PostgreSQL)                              │ │
│  │  - Agent coordination                                         │ │
│  │  - Workflow execution                                         │ │
│  └──────────────────────────┬────────────────────────────────────┘ │
│                             │                                       │
│  ┌──────────────────────────┴────────────────────────────────────┐ │
│  │              Specialized Agents (Python)                      │ │
│  │  - PM, Architect, Analyst (Planning)                          │ │
│  │  - TDD, Specialist, Reviewer (Implementation)                 │ │
│  │  - Git, QA (Delivery)                                         │ │
│  └───────────────────────────────────────────────────────────────┘ │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐ │
│  │                    Supporting Services                        │ │
│  │  ┌─────────────┬─────────────┬─────────────┬─────────────┐   │ │
│  │  │ LM Studio   │  PostgreSQL │   Redis     │  ChromaDB   │   │ │
│  │  │ (LLM API)   │  (State)    │  (Events)   │  (Vectors)  │   │ │
│  │  └─────────────┴─────────────┴─────────────┴─────────────┘   │ │
│  └───────────────────────────────────────────────────────────────┘ │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐ │
│  │                   Observability Stack                         │ │
│  │  Jaeger (Tracing) | Prometheus (Metrics) | Grafana (Viz)     │ │
│  └───────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

### Event-Driven Architecture

Everything communicates via events through Redis Streams:

**Event Types:**
```typescript
enum EventType {
  // Lifecycle
  ISSUE_CREATED = "issue.created",
  TASK_CLAIMED = "task.claimed",
  WORKTREE_CREATED = "worktree.created",
  WORKSPACE_LAUNCHED = "workspace.launched",
  WORKSPACE_DETACHED = "workspace.detached",
  
  // Planning
  PRD_GENERATED = "prd.generated",
  TRD_GENERATED = "trd.generated",
  TASKS_DECOMPOSED = "tasks.decomposed",
  
  // Implementation
  TESTS_GENERATED = "tests.generated",
  CODE_GENERATED = "code.generated",
  TESTS_PASSED = "tests.passed",
  TESTS_FAILED = "tests.failed",
  
  // Review
  REVIEW_APPROVED = "review.approved",
  REVIEW_REJECTED = "review.rejected",
  
  // Delivery
  PR_CREATED = "pr.created",
  PR_MERGED = "pr.merged",
  
  // Errors
  AGENT_FAILED = "agent.failed",
  ESCALATION_REQUIRED = "escalation.required"
}
```

**Why Event-Driven?**
- **Decoupling**: Components don't know about each other
- **Scalability**: Multiple agents work in parallel
- **Resilience**: Failures don't cascade
- **Auditability**: Complete history of every action

---

## Implementation Details

### 1. Cockpit (Go + Bubbletea)

```go
// tui/app.go
package tui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type Model struct {
    workspaceManager *workspace.Manager
    factoryClient    *client.FactoryClient
    
    // Views
    currentView View
    
    // Data
    tasks       []Task
    agents      []Agent
    worktrees   []Worktree
    
    // UI state
    selectedTask int
    width        int
    height       int
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "i", "enter":
            // Jump into workspace!
            return m, m.launchWorkspace()
        case "x":
            // Close workspace
            return m, m.closeWorkspace()
        case "a":
            // Approve PR
            return m, m.approveTask()
        }
    
    case workspaceReturnMsg:
        // User returned from workspace
        return m, m.fetchData()
    }
    
    return m, nil
}

func (m *Model) launchWorkspace() tea.Cmd {
    return func() tea.Msg {
        task := m.tasks[m.selectedTask]
        
        // Create/get workspace via Sandcastle
        ws, err := m.workspaceManager.CreateWorkspace(
            task.ID, task.IssueID, task.Branch,
        )
        if err != nil {
            return errMsg{err}
        }
        
        // Launch tmux session (blocks until detach)
        err = m.workspaceManager.LaunchWorkspace(ws)
        
        // User returned!
        return workspaceReturnMsg{taskID: task.ID}
    }
}
```

### 2. Workspace Manager (Go)

```go
// workspace/manager.go
package workspace

type Workspace struct {
    TaskID       string
    IssueID      string
    WorktreePath string
    Branch       string
    TmuxSession  string
    Active       bool
}

type Manager struct {
    worktreeRoot string
    sandcastle   *SandcastleClient
    sessions     map[string]*Workspace
}

func (m *Manager) CreateWorkspace(taskID, issueID, branch string) (*Workspace, error) {
    // Call Sandcastle API to create worktree
    result, err := m.sandcastle.CreateWorktree(branch)
    if err != nil {
        return nil, err
    }
    
    ws := &Workspace{
        TaskID:       taskID,
        IssueID:      issueID,
        WorktreePath: result.WorktreePath,
        Branch:       branch,
        TmuxSession:  fmt.Sprintf("factory-%s", issueID),
        Active:       true,
    }
    
    m.sessions[taskID] = ws
    return ws, nil
}

func (m *Manager) LaunchWorkspace(ws *Workspace) error {
    // Create tmux session
    m.createTmuxSession(ws)
    
    // Attach (blocks)
    return m.attachToSession(ws.TmuxSession)
}

func (m *Manager) createTmuxSession(ws *Workspace) error {
    sessionName := ws.TmuxSession
    
    // Create session
    exec.Command("tmux", "new-session", "-d", "-s", sessionName,
        "-c", ws.WorktreePath).Run()
    
    // Split vertically
    exec.Command("tmux", "split-window", "-h", "-t", sessionName,
        "-c", ws.WorktreePath).Run()
    
    // Left: Neovim
    exec.Command("tmux", "send-keys", "-t", sessionName+":0.0",
        "nvim .", "Enter").Run()
    
    // Right: Pi Agent
    exec.Command("tmux", "send-keys", "-t", sessionName+":0.1",
        fmt.Sprintf("pi --worktree-mode --task-id=%s", ws.TaskID),
        "Enter").Run()
    
    // Set status bar
    exec.Command("tmux", "set-option", "-t", sessionName,
        "status-left",
        fmt.Sprintf("🏭 Factory | 📂 %s | 🌿 %s", ws.IssueID, ws.Branch)).Run()
    
    return nil
}
```

### 3. Sandcastle Bridge (Node.js → Python)

```typescript
// sandcastle-service.ts
import express from 'express';
import { run, createWorktree } from "@ai-hero/sandcastle";
import { docker } from "@ai-hero/sandcastle/sandboxes/docker";

const app = express();
app.use(express.json());

app.post('/api/worktree/create', async (req, res) => {
    const { branch, taskId } = req.body;
    
    await using wt = await createWorktree({
        branchStrategy: { type: "branch", branch }
    });
    
    res.json({
        worktreePath: wt.worktreePath,
        branch: wt.branch,
        taskId
    });
});

app.post('/api/agent/run', async (req, res) => {
    const { taskId, prompt, worktreePath } = req.body;
    
    const result = await run({
        agent: customAgent,
        sandbox: docker(),
        cwd: worktreePath,
        prompt,
        maxIterations: 5
    });
    
    res.json({
        commits: result.commits,
        success: result.completionSignal === "DONE"
    });
});

app.listen(3001);
```

### 4. Python Orchestrator (Event Bus)

```python
# core/event_bus.py
class EventBus:
    def __init__(self, redis: Redis):
        self.redis = redis
        self.handlers: Dict[str, List[Callable]] = {}
    
    async def publish(self, event: Event):
        """Publish event to Redis Streams"""
        await self.redis.xadd("events", {
            "id": event.id,
            "type": event.type,
            "timestamp": str(event.timestamp),
            "source": event.source,
            "data": json.dumps(event.data),
            "trace_id": event.trace_id
        })
        
        # Also send to WebSocket for cockpit real-time updates
        await self.websocket_manager.broadcast(event)
    
    async def start_consuming(self, consumer_id: str):
        """Start consuming events"""
        while True:
            events = await self.redis.xreadgroup(
                groupname="factory",
                consumername=consumer_id,
                streams={"events": ">"},
                count=10,
                block=1000
            )
            
            for stream, messages in events:
                for message_id, data in messages:
                    await self._handle_event(message_id, data)
```

### 5. LM Studio Integration

```python
# core/llm.py
class LMStudioClient:
    def __init__(self, base_url: str = "http://localhost:1234/v1"):
        self.base_url = base_url
        self.session = aiohttp.ClientSession()
    
    async def generate(
        self,
        prompt: str,
        model_type: str = "coding",
        system: Optional[str] = None,
        temperature: float = 0.7,
        max_tokens: int = 2000
    ) -> Dict:
        """Generate completion using LM Studio"""
        
        messages = []
        if system:
            messages.append({"role": "system", "content": system})
        messages.append({"role": "user", "content": prompt})
        
        async with self.session.post(
            f"{self.base_url}/chat/completions",
            json={
                "model": MODELS[model_type].name,
                "messages": messages,
                "temperature": temperature,
                "max_tokens": max_tokens
            }
        ) as resp:
            result = await resp.json()
            return {
                "content": result["choices"][0]["message"]["content"],
                "tokens": result["usage"]["total_tokens"]
            }
```

---

## Hardware Optimization

### VRAM Budget (16GB Total)

```yaml
model_memory:
  deepseek-coder-v2-16b-q4: 9 GB
  qwen2.5-coder-14b-q4: 8 GB
  qwen2.5-coder-7b-q4: 4.5 GB
  nomic-embed-text: 0.25 GB
  
system_overhead: ~1 GB

strategies:
  balanced: "deepseek-coder-v2-16b + embeddings = 9.25GB (recommended)"
  speed: "Two qwen2.5-coder-7b instances = 9GB (parallel)"
  quality: "Single deepseek-coder-v2-16b = 9GB (sequential)"
```

### Performance Expectations

```yaml
single_task_latency:
  simple_task: 10-30 seconds
  medium_task: 1-3 minutes
  complex_task: 5-10 minutes

throughput:
  sequential: 10-15 tasks/hour
  parallel_speed_mode: 20-30 tasks/hour

daily_capacity:
  8_hour_workday_sequential: 80-120 tasks
  8_hour_workday_parallel: 160-240 tasks

token_generation_speed:
  deepseek-coder-v2-16b: 25-30 tokens/sec
  qwen2.5-coder-7b: 50-60 tokens/sec
```

---

## Implementation Roadmap

### Phase 0: Foundation (Week 1)

**Days 1-2: Infrastructure Setup**
- [ ] Install Docker + Docker Compose
- [ ] Install LM Studio + download models
- [ ] Install Go 1.22+
- [ ] Install Node.js 20+ (for Sandcastle)
- [ ] Setup PostgreSQL + Redis

**Days 3-4: Sandcastle Integration**
- [ ] Install Sandcastle (`npm install @ai-hero/sandcastle`)
- [ ] Create Sandcastle bridge service (Node.js API)
- [ ] Test worktree creation
- [ ] Test Docker sandbox provider

**Days 5-7: Python Core**
- [ ] Implement EventBus (Redis Streams)
- [ ] Implement LMStudioClient
- [ ] Implement VectorStore (ChromaDB)
- [ ] Test end-to-end event flow

### Phase 1: Cockpit TUI (Week 2)

**Days 1-3: Basic Cockpit**
- [ ] Setup Go project with Bubbletea
- [ ] Implement dashboard view
- [ ] Implement task list view
- [ ] Add keyboard navigation

**Days 4-5: Workspace Integration**
- [ ] Implement workspace manager (Go)
- [ ] Create tmux session launcher
- [ ] Test Neovim + Pi split view
- [ ] Add detach/reattach logic

**Days 6-7: Real-Time Updates**
- [ ] WebSocket client in Go
- [ ] Connect to Python event bus
- [ ] Update UI on events
- [ ] Test live monitoring

### Phase 2: Agents (Week 3-4)

**Week 3: Planning Agents**
- [ ] Implement BaseAgent class
- [ ] Implement PMAgent
- [ ] Implement ArchitectAgent
- [ ] Implement AnalystAgent
- [ ] Test planning workflow

**Week 4: Implementation Agents**
- [ ] Implement TDDAgent
- [ ] Implement SpecialistAgent (Python)
- [ ] Implement ReviewerAgent
- [ ] Implement Sandbox executor
- [ ] Test implementation workflow

### Phase 3: Integration & Polish (Week 5-6)

**Week 5: End-to-End Testing**
- [ ] Full workflow: Issue → Planning → Implementation → PR
- [ ] Test workspace persistence
- [ ] Test multiple concurrent tasks
- [ ] Performance optimization

**Week 6: UI/UX Polish**
- [ ] Add more cockpit views (metrics, logs, settings)
- [ ] Improve styling with Lipgloss
- [ ] Add help screens
- [ ] Create documentation
- [ ] Record demo video

### Total Timeline: ~6 weeks to working system

---

## Code Reference

### Project Structure

```
dev-factory/
├── cockpit/                   # Go TUI
│   ├── main.go
│   ├── go.mod
│   ├── tui/
│   │   ├── app.go            # Main Bubbletea app
│   │   ├── dashboard.go      # Dashboard view
│   │   ├── tasks.go          # Task list
│   │   ├── agents.go         # Agent monitor
│   │   └── styles.go         # Lipgloss styles
│   ├── workspace/
│   │   ├── manager.go        # Workspace manager
│   │   └── tmux.go           # Tmux integration
│   └── client/
│       ├── factory.go        # API client
│       └── websocket.go      # Real-time updates
│
├── sandcastle-bridge/         # Node.js service
│   ├── package.json
│   ├── server.ts
│   └── .sandcastle/
│       ├── prompt.md
│       └── main.ts
│
├── orchestrator/              # Python core
│   ├── core/
│   │   ├── event_bus.py
│   │   ├── llm.py
│   │   ├── vector_store.py
│   │   └── sandbox.py
│   ├── agents/
│   │   ├── base_agent.py
│   │   ├── pm_agent.py
│   │   ├── specialist_agent.py
│   │   └── reviewer_agent.py
│   ├── workflows/
│   │   ├── planning.py
│   │   └── implementation.py
│   └── api/
│       ├── main.py           # FastAPI
│       └── websocket.py      # WebSocket handler
│
├── config/
│   ├── models.yaml
│   ├── agents.yaml
│   ├── workspace-layouts.yaml
│   └── prometheus.yml
│
├── scripts/
│   ├── setup.sh
│   └── build-cockpit.sh
│
└── docker-compose.yml
```

---

## Next Steps

### Immediate (This Week)

1. **Setup infrastructure**
   ```bash
   # Install dependencies
   sudo apt install docker docker-compose tmux neovim golang nodejs npm
   
   # Install LM Studio
   # Download from https://lmstudio.ai
   
   # Install Sandcastle
   npm install -g @ai-hero/sandcastle
   ```

2. **Initialize projects**
   ```bash
   cd ~/Projects
   mkdir dev-factory && cd dev-factory
   
   # Go cockpit
   mkdir cockpit && cd cockpit
   go mod init github.com/yourusername/factory-cockpit
   go get github.com/charmbracelet/bubbletea
   go get github.com/charmbracelet/lipgloss
   
   # Python orchestrator
   cd ..
   mkdir orchestrator && cd orchestrator
   python -m venv venv
   source venv/bin/activate
   pip install fastapi uvicorn asyncpg redis aiohttp
   
   # Sandcastle bridge
   cd ..
   mkdir sandcastle-bridge && cd sandcastle-bridge
   npm init -y
   npm install @ai-hero/sandcastle express
   ```

3. **Test basic integration**
   - Start LM Studio server
   - Create test worktree with Sandcastle
   - Launch tmux session manually
   - Verify Neovim + Pi work together

### Short Term (Weeks 2-3)

1. **Build cockpit MVP**
   - Dashboard view
   - Task list
   - Workspace launcher

2. **Implement core agents**
   - TDD Agent
   - Specialist Agent
   - Basic workflow

3. **Test end-to-end**
   - Create task → Jump into workspace → Code → Return

### Medium Term (Weeks 4-6)

1. **Complete agent suite**
2. **Add all cockpit views**
3. **Polish UX**
4. **Write documentation**
5. **Create demo**

---

## Key Insights

### Why This Architecture?

**1. Terminal-First Approach**
- Developers live in the terminal
- No context switching to browser
- SSH-friendly
- Keyboard-driven workflow
- Fast and responsive

**2. Integrated Workspace**
- Seamless Neovim + AI experience
- No switching between windows
- Task context is preserved
- Natural development flow

**3. Event-Driven Core**
- Scalable and resilient
- Easy to add new agents
- Full audit trail
- Decoupled components

**4. Local-First**
- Zero API costs
- Complete privacy
- Full control
- Works offline

**5. Best-of-Breed Tools**
- Bubbletea for beautiful TUIs
- Sandcastle for worktree management
- LM Studio for local LLMs
- Tmux for multiplexing

### Comparison to Alternatives

**vs. Web-Based IDEs (VS Code, Cursor)**
- ✅ Lives in terminal (no Electron)
- ✅ SSH-friendly
- ✅ More customizable
- ❌ Requires terminal knowledge

**vs. Traditional Dev Workflow**
- ✅ AI-native experience
- ✅ Task isolation
- ✅ Automated workflows
- ✅ Beautiful monitoring

**vs. Cloud AI Tools (GitHub Copilot, Cursor)**
- ✅ 100% local and free
- ✅ No usage limits
- ✅ Complete privacy
- ❌ Requires hardware

---

## Conclusion

This is a **production-ready architecture** for a beautiful, terminal-native, AI-powered development factory with:

1. **Zero cost** (all free and open source)
2. **Beautiful UI** (Go + Bubbletea cockpit)
3. **Integrated workspace** (Neovim + Pi in tmux)
4. **Worktree isolation** (Sandcastle)
5. **Local LLMs** (LM Studio)
6. **Event-driven** (Python + Redis Streams)
7. **Observable** (Prometheus + Grafana)

The cockpit provides a **mission control experience** while the integrated workspace provides a **seamless IDE experience**. Together, they create a uniquely powerful and enjoyable development environment.

**Status**: Ready for implementation  
**Owner**: m00nk0d3  
**Last Updated**: 2026-07-12  
**Next Action**: Run setup script and build cockpit MVP
