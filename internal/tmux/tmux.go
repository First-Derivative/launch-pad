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
//   - Window "code": pane 0 (90%) runs nvim, pane 1 (10%) runs oc
//   - Window "tools": pane 0 (30%) runs lazygit, pane 1 (70%) empty
func CreateSession(name, dir string) error {
	commands := []struct {
		args []string
	}{
		// Create detached session with "code" window
		{[]string{"new-session", "-d", "-s", name, "-n", "code", "-c", dir}},

		// Split "code" window horizontally, new pane is 10%
		{[]string{"split-window", "-h", "-t", name + ":code", "-p", "10", "-c", dir}},

		// Send commands to "code" window panes
		{[]string{"send-keys", "-t", name + ":code.0", "nvim .", "Enter"}},
		{[]string{"send-keys", "-t", name + ":code.1", "oc", "Enter"}},

		// Create "tools" window
		{[]string{"new-window", "-t", name, "-n", "tools", "-c", dir}},

		// Split "tools" window horizontally, new pane is 70%
		{[]string{"split-window", "-h", "-t", name + ":tools", "-p", "70", "-c", dir}},

		// Send lazygit command to "tools" window left pane
		{[]string{"send-keys", "-t", name + ":tools.0", `eval "$(ssh-agent -s)" && ssh-add ~/ssh-keys/opus && lazygit`, "Enter"}},

		// Select "code" window and left pane
		{[]string{"select-window", "-t", name + ":code"}},
		{[]string{"select-pane", "-t", name + ":code.0"}},
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
//   - Window "code": pane 0 (80%) runs nvim, pane 1 (20%) runs oc
//   - Window "tools": 4 panes with lazygit, api-gateway, web-app, and docker compose
func CreatePLPSession() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	name := "plp"
	dir := homeDir + "/dev/plp-mono"

	commands := []struct {
		args []string
	}{
		// Create detached session with "code" window
		{[]string{"new-session", "-d", "-s", name, "-n", "code", "-c", dir}},

		// Split "code" window horizontally, new pane is 20%
		{[]string{"split-window", "-h", "-t", name + ":code", "-l", "20%", "-c", dir}},

		// Send commands to "code" window panes
		{[]string{"send-keys", "-t", name + ":code.0", "nvim .", "Enter"}},
		{[]string{"send-keys", "-t", name + ":code.1", "oc", "Enter"}},

		// Create "tools" window
		{[]string{"new-window", "-t", name, "-n", "tools", "-c", dir}},

		// Split "tools" window horizontally (30% left, 70% right)
		{[]string{"split-window", "-h", "-t", name + ":tools", "-l", "70%", "-c", dir}},

		// Split the right pane (pane 1) vertically into top and bottom (30% top, 70% bottom)
		{[]string{"split-window", "-v", "-t", name + ":tools.1", "-l", "70%", "-c", dir}},

		// Split the top right pane (pane 1) horizontally into left and right (50/50)
		{[]string{"split-window", "-h", "-t", name + ":tools.1", "-c", dir}},

		// Final layout:
		// - Pane 0: Left (30%) - lazygit
		// - Pane 1: Top left of right area
		// - Pane 2: Top right of right area
		// - Pane 3: Bottom (70% height)

		// Send commands to "tools" window panes
		{[]string{"send-keys", "-t", name + ":tools.0", `eval "$(ssh-agent -s)" && ssh-add ~/ssh-keys/opus && lazygit`, "Enter"}},
		{[]string{"send-keys", "-t", name + ":tools.1", "yarn workspace api-gateway-lambdas run dev", "Enter"}},
		{[]string{"send-keys", "-t", name + ":tools.2", "yarn workspace web-app run dev", "Enter"}},
		{[]string{"send-keys", "-t", name + ":tools.3", "docker compose up -d", "Enter"}},

		// Select "code" window and left pane
		{[]string{"select-window", "-t", name + ":code"}},
		{[]string{"select-pane", "-t", name + ":code.0"}},
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
