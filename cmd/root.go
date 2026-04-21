package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ashhatz/launch-pad/internal/tmux"
	"github.com/spf13/cobra"
)

var (
	createName  string // -t flag: override session name
	attachName  string // -a flag: attach to named session
	profileName string // -p flag: load a named profile configuration
)

var rootCmd = &cobra.Command{
	Use:   "launch [path]",
	Short: "A tmux session launcher",
	Long: `Launch creates and manages tmux development sessions.

By default, creates a new session and attaches to it.
Session name is derived from the directory basename, or use -t to override.

Examples:
  launch                          # Create+attach, name from CWD
  launch ~/dev/foo                # Create+attach, name="foo", dir=~/dev/foo
  launch -t myproj                # Create+attach, name="myproj", dir=CWD
  launch -t myproj ~/dev/foo      # Create+attach, name="myproj", dir=~/dev/foo
  launch -a mysession             # Attach to existing "mysession"
  launch -p plp                   # PLP profile, default name+dir
  launch -p plp ~/dev/plp-wt      # PLP profile, name from path basename
  launch -p plp -t feat ~/dev/plp-wt  # PLP profile, custom name+dir
  launch -p ai                        # AI profile, oc + claude side by side
  launch -p ai -t my-ai ~/dev/foo     # AI profile, custom name+dir`,
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

func init() {
	rootCmd.Flags().StringVarP(&attachName, "attach", "a", "", "Attach to a named session (ignores path and -t)")
	rootCmd.Flags().StringVarP(&createName, "create", "t", "", "Override session name when creating")
	rootCmd.Flags().StringVarP(&profileName, "profile", "p", "", "Load a named profile configuration (combines with -t and path)")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	// PROFILE MODE: -p takes highest priority, ignores all other flags
	if cmd.Flags().Changed("profile") {
		if profileName == "" {
			return fmt.Errorf("-p requires a profile name")
		}
		switch strings.ToLower(profileName) {
		case "plp":
			// Resolve directory: positional arg > profile default
			var dir string
			if len(args) > 0 {
				d, err := filepath.Abs(args[0])
				if err != nil {
					return fmt.Errorf("failed to resolve path '%s': %w", args[0], err)
				}
				dir = d
			} else {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home directory: %w", err)
				}
				dir = homeDir + "/dev/plp-mono"
			}

			// Resolve session name: -t flag > path basename > profile default
			var sessionName string
			if cmd.Flags().Changed("create") {
				if createName == "" {
					return fmt.Errorf("-t requires a session name")
				}
				sessionName = createName
			} else if len(args) > 0 {
				sessionName = filepath.Base(dir)
			} else {
				sessionName = "plp"
			}

			if tmux.HasSession(sessionName) {
				return fmt.Errorf("session '%s' already exists. Use -a %s to attach", sessionName, sessionName)
			}
			if err := tmux.CreatePLPSession(sessionName, dir); err != nil {
				return err
			}
			return tmux.AttachSession(sessionName)
		case "ai":
			// Resolve directory: positional arg > CWD
			var dir string
			if len(args) > 0 {
				d, err := filepath.Abs(args[0])
				if err != nil {
					return fmt.Errorf("failed to resolve path '%s': %w", args[0], err)
				}
				dir = d
			} else {
				d, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get working directory: %w", err)
				}
				dir = d
			}

			// Resolve session name: -t flag > path basename > profile default
			var sessionName string
			if cmd.Flags().Changed("create") {
				if createName == "" {
					return fmt.Errorf("-t requires a session name")
				}
				sessionName = createName
			} else if len(args) > 0 {
				sessionName = filepath.Base(dir)
			} else {
				sessionName = "ai"
			}

			if tmux.HasSession(sessionName) {
				return fmt.Errorf("session '%s' already exists. Use -a %s to attach", sessionName, sessionName)
			}
			if err := tmux.CreateAISession(sessionName, dir); err != nil {
				return err
			}
			return tmux.AttachSession(sessionName)
		default:
			return fmt.Errorf("unknown profile '%s'", profileName)
		}
	}

	// ATTACH MODE: -a takes priority, ignores path and -t
	if cmd.Flags().Changed("attach") {
		if attachName == "" {
			return fmt.Errorf("-a requires a session name")
		}
		if !tmux.HasSession(attachName) {
			return fmt.Errorf("session '%s' does not exist", attachName)
		}
		return tmux.AttachSession(attachName)
	}

	// CREATE MODE
	// 1. Resolve path (default to current directory)
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path '%s': %w", path, err)
	}

	// 2. Determine session name
	// Priority: -t value > path basename
	var sessionName string
	if createName != "" {
		sessionName = createName
	} else {
		sessionName = filepath.Base(absPath)
	}

	// 3. Create session (error if already exists)
	if tmux.HasSession(sessionName) {
		return fmt.Errorf("session '%s' already exists. Use -a %s to attach", sessionName, sessionName)
	}
	if err := tmux.CreateSession(sessionName, absPath); err != nil {
		return err
	}

	// 4. Attach to the created session
	return tmux.AttachSession(sessionName)
}
