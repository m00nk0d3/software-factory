#!/bin/bash

# Kill existing session if it exists
tmux kill-session -t factory-issue-1 2>/dev/null || true

# Create new session
tmux new-session -d -s factory-issue-1 -c ~/worktrees/issue-1

# Split window horizontally
tmux split-window -h -t factory-issue-1:0 -c ~/worktrees/issue-1

# Left pane: Neovim
tmux send-keys -t factory-issue-1:0.0 'nvim .' C-m

# Right pane: Pi Agent welcome screen
tmux send-keys -t factory-issue-1:0.1 'clear' C-m
sleep 0.3
tmux send-keys -t factory-issue-1:0.1 'cat << "EOF"
╔════════════════════════════════════════════════════════╗
║  🤖 Pi Agent - Worktree Mode                         ║
╠════════════════════════════════════════════════════════╣
║  Issue: #1 - Infrastructure Setup                    ║
║  Branch: feature/issue-1-infrastructure-setup         ║
║  Worktree: ~/worktrees/issue-1                        ║
╠════════════════════════════════════════════════════════╣
║  Tasks:                                               ║
║  [ ] Run ./scripts/setup.sh                           ║
║  [ ] Check Docker services                            ║
║  [ ] Verify LM Studio                                 ║
║  [ ] Test infrastructure                              ║
╠════════════════════════════════════════════════════════╣
║  Ready to assist! What do you need help with?        ║
╚════════════════════════════════════════════════════════╝

EOF
' C-m

# Set pane titles
tmux select-pane -t factory-issue-1:0.0 -T 'NEOVIM'
tmux select-pane -t factory-issue-1:0.1 -T 'PI AGENT'

# Configure status bar
tmux set-option -t factory-issue-1 status-left '🏭 Factory | 📂 issue-1 | 🌿 feature/issue-1 '
tmux set-option -t factory-issue-1 status-left-length 60
tmux set-option -t factory-issue-1 status-right '[Ctrl+B D] Exit to Cockpit'

# Focus on Neovim pane
tmux select-pane -t factory-issue-1:0.0

echo ""
echo "✅ Integrated workspace launched!"
echo ""
echo "Layout:"
echo "  Left:  Neovim (editing files)"
echo "  Right: Pi Agent (your AI assistant)"
echo ""
echo "Controls:"
echo "  Ctrl+B ←/→  - Switch between panes"
echo "  Ctrl+B D    - Detach (return to shell)"
echo "  Ctrl+B Z    - Zoom current pane"
echo ""
echo "Attaching in 2 seconds..."
sleep 2

# Attach to session
tmux attach-session -t factory-issue-1
