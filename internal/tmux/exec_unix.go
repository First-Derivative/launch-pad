//go:build unix

package tmux

import "syscall"

// execSyscall replaces the current process with the given command.
// This is used for tmux attach to work properly.
func execSyscall(path string, args []string, env []string) error {
	return syscall.Exec(path, args, env)
}
