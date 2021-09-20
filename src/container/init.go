/*
 * @Author: dejavudwh
 * @Date: 2021-09-06 17:10:24
 * @LastEditTime: 2021-09-20 12:06:59
 */
package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}

	setUpMount()

	log.Info(cmdArray)
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

/**
 * @description: init mount point
 * @param {*}
 * @return {*}
 */
func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
	}

	log.Infof("Current location is %s", pwd)
	// switch root
	pivotRoot(pwd)

	/*
		* mount proc to view process resources later (in mount namespace)
		- MS_NOEXEC: Do not allow other programs to run in this filesystem
		- MS_NOSUID: When running the program, set-user-ID or set-group-ID is not allowed
		- MS_NODEV: default parameter
	*/
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	// mount tmpfs, tmpfs is a memory-based file system
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

/**
 * @description: mount root file system
 * @param {*}: dir path to be replaced as root
 * @return {*}
 */

func pivotRoot(root string) error {
	/*
		* pivot_root is a ways to switch old root to new root
		* bind mount is just a way to replace the same content with a mount point
		- MS_BIND: Create a bind mount
		- MS_REC: Used in conjunction with MS_BIND to create a recursive bind mount
	*/
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}

	// create roofs/.pivot_root to store old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	// pivot_root to new rootfs(Now the old old_root is mounted in rootfs/.pivot_root)
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	// modify the current working directory to the root directory
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// unmount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	// remove dir
	return os.Remove(pivotDir)
}
