# Software Factory 🏭

> A fully autonomous, local-first software development factory powered by AI agents

## What is this?

An end-to-end development automation system that:
- 🤖 Uses specialized AI agents for different stages of development
- 🏠 Runs 100% locally with LM Studio (no API costs)
- 🎮 Features a beautiful Go TUI cockpit for monitoring and control
- 💻 Provides integrated workspaces (Neovim + Pi Agent in tmux)
- 🔄 Manages isolated git worktrees via Sandcastle
- ⚡ Event-driven architecture for scalability

## Architecture

```
┌──────────────────────────────────────────┐
│    Factory Cockpit (Go + Bubbletea)     │
│    Beautiful terminal dashboard          │
└────────────┬─────────────────────────────┘
             │
             ▼
┌──────────────────────────────────────────┐
│   Integrated Workspace (Tmux)            │
│   Neovim  |  Pi Agent                    │
└────────────┬─────────────────────────────┘
             │
             ▼
┌──────────────────────────────────────────┐
│   Python Orchestrator (Event-Driven)     │
│   Specialized Agents + Workflows         │
└────────────┬─────────────────────────────┘
             │
             ▼
┌──────────────────────────────────────────┐
│   Sandcastle (Worktree Management)       │
│   Git worktree isolation per task        │
└──────────────────────────────────────────┘
```

## Tech Stack

- **Cockpit**: Go + Bubbletea (TUI framework)
- **Orchestrator**: Python + asyncio + Redis Streams
- **Worktrees**: Sandcastle (TypeScript)
- **LLM**: LM Studio (deepseek-coder-v2-16b, qwen2.5-coder)
- **Database**: PostgreSQL + Redis
- **Vectors**: ChromaDB
- **Observability**: Prometheus + Grafana + Jaeger

## Hardware Requirements

- **Minimum**: 8 cores, 32GB RAM, GTX 1660 (6GB VRAM)
- **Recommended**: 16 cores, 64GB RAM, RTX 3090/4090 (24GB VRAM)
- **Current Setup**: 32GB RAM, RTX 5070 Ti (16GB VRAM) ✅

## Features

### 🎮 Factory Cockpit
- Real-time dashboard with task queue, active agents, and metrics
- Jump into any task's workspace with a single keypress
- Live log streaming and approval queue
- Beautiful terminal UI that feels like a spaceship control panel

### 💻 Integrated Workspace
- Neovim + Pi Agent side-by-side in tmux
- Task context automatically loaded
- Persistent sessions (detach/reattach anytime)
- Multiple layout options

### 🤖 Specialized Agents
- **Planning**: PM, Architect, Analyst
- **Implementation**: TDD, Specialist, Reviewer
- **Delivery**: Git, QA, PR Manager

### 🔄 Event-Driven
- All components communicate via Redis Streams
- Parallel agent execution
- Full audit trail
- Easy to extend

## Project Structure

```
software-factory/
├── cockpit/                # Go TUI cockpit
├── orchestrator/           # Python event-driven core
├── sandcastle-bridge/      # Node.js Sandcastle API
├── config/                 # Configuration files
├── scripts/               # Setup and utility scripts
└── docker-compose.yml     # Infrastructure services
```

## Quick Start

```bash
# Clone the repo
git clone https://github.com/yourusername/software-factory.git
cd software-factory

# Run setup script
./scripts/setup.sh

# Start the cockpit
./cockpit/factory-cockpit
```

## Status

🚧 **Under Active Development** 🚧

Current phase: Initial implementation

See [implementation guide](https://github.com/yourusername/software-factory/blob/main/docs/implementation.md) for detailed architecture and roadmap.

## License

MIT

## Author

Built with 🔥 by m00nk0d3
