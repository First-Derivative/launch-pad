package tmux

import (
	"fmt"
	"os"
	"os/exec"
)

// HasSession checks if a tmux session with the given name exists.
func HasSession(name string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", name)
	err := cmd.Run()
	return err == nil
}

// CreateSession creates a new detached tmux session with the standard layout.
// Layout:
//   - Window "code": pane 0 (80%) runs nvim, pane 1 (20%) runs lazygit
//   - Window "ai": pane 0 (50%) runs oc, pane 1 (50%) runs claude
//   - Window "bash": pane 0 (50%) empty, pane 1 (50%) empty
func CreateSession(name, dir string) error {
	commands := []struct {
		args []string
	}{
		// Create detached session with "code" window
		{[]string{"new-session", "-d", "-s", name, "-n", "code", "-c", dir}},

		// Split "code" window horizontally: pane 0 = 80% nvim, pane 1 = 20% lazygit
		{[]string{"split-window", "-h", "-t", name + ":code", "-p", "20", "-c", dir}},

		// Send commands to "code" window panes
		{[]string{"send-keys", "-t", name + ":code.0", "nvim .", "Enter"}},
		{[]string{"send-keys", "-t", name + ":code.1", `eval "$(ssh-agent -s)" && ssh-add ~/ssh-keys/opus && lazygit`, "Enter"}},

		// Create "ai" window
		{[]string{"new-window", "-t", name, "-n", "ai", "-c", dir}},

		// Split "ai" window horizontally 50/50
		{[]string{"split-window", "-h", "-t", name + ":ai", "-p", "50", "-c", dir}},

		// Send commands to "ai" window panes
		{[]string{"send-keys", "-t", name + ":ai.0", "oc", "Enter"}},
		{[]string{"send-keys", "-t", name + ":ai.1", "claude", "Enter"}},

		// Create "bash" window
		{[]string{"new-window", "-t", name, "-n", "bash", "-c", dir}},

		// Split "bash" window horizontally 50/50
		{[]string{"split-window", "-h", "-t", name + ":bash", "-p", "50", "-c", dir}},

		// Select "ai" window
		{[]string{"select-window", "-t", name + ":ai"}},
	}

	for _, c := range commands {
		cmd := exec.Command("tmux", c.args...)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("tmux %v failed: %w", c.args[0], err)
		}
	}

	return nil
}

// CreatePLPSession creates a new detached tmux session with the PLP project layout.
// Layout:
//   - Window "command-center": single pane runs command-center
//   - Window "code": pane 0 (80%) runs nvim, pane 1 (20%) runs lazygit
//   - Window "ai": pane 0 (50%) runs oc, pane 1 (50%) runs claude
//   - Window "bash": 3 panes — empty bash (top-left), web-app (top-right), docker compose (bottom)
func CreatePLPSession(name, dir string) error {
	commands := []struct {
		args []string
	}{
		// Create detached session with "command-center" window
		{[]string{"new-session", "-d", "-s", name, "-n", "command-center", "-c", dir}},

		// Send command to "command-center" window
		{[]string{"send-keys", "-t", name + ":command-center", "command-center", "Enter"}},

		// Create "code" window
		{[]string{"new-window", "-t", name, "-n", "code", "-c", dir}},

		// Split "code" window horizontally: pane 0 = 80% nvim, pane 1 = 20% lazygit
		{[]string{"split-window", "-h", "-t", name + ":code", "-l", "20%", "-c", dir}},

		// Send commands to "code" window panes
		{[]string{"send-keys", "-t", name + ":code.0", "nvim .", "Enter"}},
		{[]string{"send-keys", "-t", name + ":code.1", `eval "$(ssh-agent -s)" && ssh-add ~/ssh-keys/opus && lazygit`, "Enter"}},

		// Create "ai" window
		{[]string{"new-window", "-t", name, "-n", "ai", "-c", dir}},

		// Split "ai" window horizontally 50/50
		{[]string{"split-window", "-h", "-t", name + ":ai", "-l", "50%", "-c", dir}},

		// Send commands to "ai" window panes
		{[]string{"send-keys", "-t", name + ":ai.0", "oc", "Enter"}},
		{[]string{"send-keys", "-t", name + ":ai.1", "claude", "Enter"}},

		// Create "bash" window
		{[]string{"new-window", "-t", name, "-n", "bash", "-c", dir}},

		// Split "bash" window vertically: top 30%, bottom 70%
		{[]string{"split-window", "-v", "-t", name + ":bash", "-l", "70%", "-c", dir}},

		// Split top pane horizontally 50/50
		{[]string{"split-window", "-h", "-t", name + ":bash.0", "-c", dir}},

		// Final layout:
		// - Pane 0: Top left - empty bash
		// - Pane 1: Top right - web-app
		// - Pane 2: Bottom (70% height) - docker compose

		// Send commands to "bash" window panes
		{[]string{"send-keys", "-t", name + ":bash.1", "yarn workspace web-app run dev", "Enter"}},
		{[]string{"send-keys", "-t", name + ":bash.2", "docker compose up -d", "Enter"}},

		// Select "ai" window
		{[]string{"select-window", "-t", name + ":ai"}},
	}

	for _, c := range commands {
		cmd := exec.Command("tmux", c.args...)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("tmux %v failed: %w", c.args[0], err)
		}
	}

	return nil
}

// CreateAISession creates a new detached tmux session with two panes for AI tools.
// Layout:
//   - Window "ai": pane 0 (50%) runs oc, pane 1 (50%) runs claude
func CreateAISession(name, dir string) error {
	commands := []struct {
		args []string
	}{
		// Create detached session with "ai" window
		{[]string{"new-session", "-d", "-s", name, "-n", "ai", "-c", dir}},

		// Split "ai" window horizontally 50/50
		{[]string{"split-window", "-h", "-t", name + ":ai", "-l", "50%", "-c", dir}},

		// Send commands to "ai" window panes
		{[]string{"send-keys", "-t", name + ":ai.0", "oc", "Enter"}},
		{[]string{"send-keys", "-t", name + ":ai.1", "claude", "Enter"}},
	}

	for _, c := range commands {
		cmd := exec.Command("tmux", c.args...)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("tmux %v failed: %w", c.args[0], err)
		}
	}

	return nil
}

// AttachSession attaches to an existing tmux session.
// This replaces the current process with tmux.
func AttachSession(name string) error {
	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf("tmux not found: %w", err)
	}

	// Use syscall.Exec to replace the current process
	// This is necessary for tmux attach to work properly
	return execSyscall(tmuxPath, []string{"tmux", "attach-session", "-t", name}, os.Environ())
}
