/*
 * @Author: dejavudwh
 * @Date: 2021-09-06 14:59:22
 * @LastEditTime: 2021-09-19 17:32:20
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
	mntURL := "/home/mnt/"
	rootURL := "/home/"
	NewWorkSpace(rootURL, mntURL)
	// image has been mounted to mntURL
	cmd.Dir = mntURL

	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	return read, write, nil
}

/**
 * @description: create a overlayfs filesystem as container root workspace
 * @param {string} rootURL: image path
 * @param {string} mntURL: mount path
 * @return {*}
 */
func NewWorkSpace(rootURL string, mntURL string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)
}

func CreateReadOnlyLayer(rootURL string) {
	busyboxURL := rootURL + "busybox/"
	busyboxTarURL := rootURL + "busybox.tar"

	exist, err := PathExists(busyboxURL)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", busyboxURL, err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s error. %v", busyboxURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error %v", busyboxURL, err)
		}
	}
}

func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"
	if err := os.Mkdir(writeURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeURL, err)
	}
}

func CreateMountPoint(rootURL string, mntURL string) {
	if err := os.Mkdir(mntURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mntURL, err)
	}

	dirs := "lowerdir=" + rootURL + "busybox" + "," + "upperdir=" + rootURL + "writeLayer" + "," + "workdir=" + rootURL + "tmpwork"
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
