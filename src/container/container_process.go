package container

import (
	"os"
	"os/exec"
	"syscall"
)

/*
* Set up the environment and create a new command to create a new process
 */
func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	// Copy its own process as the initialization of the new process
	// will call initCommand
	cmd := exec.Command("/proc/self/exe", args...)
	// namespace argument of new process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	// Process IO redirection
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}
