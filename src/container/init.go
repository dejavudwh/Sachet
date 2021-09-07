/*
 * @Author: dejavudwh
 * @Date: 2021-09-06 17:10:24
 * @LastEditTime: 2021-09-07 18:38:39
 */
package container

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

/**
 * @description: Container is created, this function starts to initialize container
 * @return {*}
 */
func RunContainerInitProcess() error {
	cmdArray := readUserCommand()

	/*
		* mount proc to view process resources later (in mount namespace)
		- MS_NOEXEC: Do not allow other programs to run in this filesystem
		- MS_NOSUID: When running the program, set-user-ID or set-group-ID is not allowed
		- MS_NODEV: default parameter
	*/
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}
	log.Infof("Find path %s", path)
	// syscall.Exec will overwrite the init process, and the user process becomes the first process
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}

	return nil
}

func readUserCommand() []string {
	// fd is 3, first file of child process
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)

	return strings.Split(msgStr, " ")
}
