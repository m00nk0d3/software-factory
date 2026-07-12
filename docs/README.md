# Software Factory Documentation

Welcome to the Software Factory documentation! This directory contains all technical documentation, architecture details, and implementation guides.

## 📚 Documentation Index

### Core Documentation

- **[Implementation Guide](./implementation.md)** - Complete architecture and implementation roadmap
  - Executive summary
  - Tech stack details
  - Cockpit UI design
  - Integrated workspace architecture
  - Event-driven system design
  - Hardware optimization
  - Complete code examples
  - 6-week implementation roadmap

## 🎯 Quick Navigation

### For New Contributors
Start with the [Implementation Guide](./implementation.md) - it covers everything you need to know about the system architecture and design decisions.

### For Users
- **Getting Started**: See main [README](../README.md)
- **Hardware Requirements**: [Implementation Guide - Hardware Optimization](./implementation.md#hardware-optimization)
- **Quick Start**: [Implementation Guide - Next Steps](./implementation.md#next-steps)

### For Developers
- **Architecture**: [Implementation Guide - Core Architecture](./implementation.md#core-architecture)
- **Tech Stack**: [Implementation Guide - Final Tech Stack](./implementation.md#final-tech-stack)
- **Code Reference**: [Implementation Guide - Code Reference](./implementation.md#code-reference)
- **Roadmap**: [Implementation Guide - Implementation Roadmap](./implementation.md#implementation-roadmap)

## 📖 What's Documented

### System Architecture
- Factory Cockpit (Go + Bubbletea TUI)
- Integrated Workspace (Neovim + Pi Agent in tmux)
- Event-Driven Orchestration (Python + Redis Streams)
- Worktree Management (Sandcastle)
- LLM Integration (LM Studio)

### Key Components
- **Cockpit**: Beautiful terminal dashboard with real-time monitoring
- **Workspace**: Seamless Neovim + AI assistant experience
- **Agents**: Specialized AI agents for each development phase
- **Events**: Async message passing for scalability
- **Worktrees**: Isolated git worktrees per task

### Implementation Details
- Complete code examples for all components
- Database schemas
- API specifications
- Configuration formats
- Deployment architecture

## 🚀 Getting Started with Docs

1. Read the [Implementation Guide](./implementation.md) from top to bottom for full context
2. Focus on sections relevant to your interest:
   - **Building the Cockpit**: See "The Factory Cockpit" section
   - **Understanding Agents**: See "Core Architecture" section
   - **Setting Up Infrastructure**: See "Implementation Details" section
3. Follow the roadmap in "Implementation Roadmap" section

## 📝 Contributing to Docs

If you're adding new documentation:
1. Keep the same structure and formatting
2. Add code examples where relevant
3. Update this README index
4. Use clear section headers for navigation

## 🔗 External Resources

- **Bubbletea**: https://github.com/charmbracelet/bubbletea
- **Sandcastle**: https://github.com/mattpocock/sandcastle
- **LM Studio**: https://lmstudio.ai
- **Redis Streams**: https://redis.io/docs/data-types/streams/

## 💡 Philosophy

This documentation follows these principles:
- **Complete**: Everything needed to build the system from scratch
- **Practical**: Real code examples, not just theory
- **Clear**: Simple language, well-structured
- **Visual**: Diagrams and ASCII art where helpful
- **Actionable**: Step-by-step guides and roadmaps

---

**Last Updated**: 2026-07-12  
**Status**: Under active development  
**Maintainer**: m00nk0d3
