/*
 * @Author: dejavudwh
 * @Date: 2021-09-06 14:59:22
 * @LastEditTime: 2021-09-06 18:28:39
 */
package container

import (
	"os"
	"os/exec"
	"syscall"
)

/**
 * @description: Set up the environment and create a new command to create a new process
 * @param {bool} tty: Process IO redirection
 * @param {string} command: user first command and init command
 * @return {*}
 */
func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	// Copy its own process as the initialization of the new process
	// Will call initCommand
	cmd := exec.Command("/proc/self/exe", args...)
	// namespace argument of new process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}
