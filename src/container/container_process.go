/*
 * @Author: dejavudwh
 * @Date: 2021-09-06 14:59:22
 * @LastEditTime: 2021-09-07 15:50:58
 */
package container

import (
	"os"
	"os/exec"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

/**
 * @description: Set up the environment and create a new command to create a new process
 * @param {bool} tty: Process IO redirection
 * @param {string} command: user first command and init command
 * @return {*}
 */
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}
	// Copy its own process as the initialization of the new process (child process)
	// Will call initCommand
	cmd := exec.Command("/proc/self/exe", "init")
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

	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	return read, write, nil
}
