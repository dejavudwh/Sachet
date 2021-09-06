/*
 * @Author: dejavudwh
 * @Date: 2021-09-06 17:10:24
 * @LastEditTime: 2021-09-06 17:56:23
 */
package container

import (
	"os"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

/**
 * @description: Container is created, this function starts to initialize container
 * @param {string} command: first command
 * @param {[]string} args: command arugument
 * @return {*}
 */
func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command %s", command)

	/*
		* mount proc to view process resources later (in mount namespace)
		- MS_NOEXEC: Do not allow other programs to run in this filesystem
		- MS_NOSUID: When running the program, set-user-ID or set-group-ID is not allowed
		- MS_NODEV: default parameter
	*/
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	// Will overwrite the init process, and the user process becomes the first process
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}

	return nil
}
