/*
 * @Author: dejavudwh
 * @Date: 2021-09-06 14:08:11
 * @LastEditTime: 2021-09-19 17:40:02
 */
package main

import (
	"Scachet/src/cgroups"
	"Scachet/src/cgroups/subsystems"
	"Scachet/src/container"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

/**
 * @description: Create a new process through the created command
 * @param {bool} tty: IO redirection
 * @param {string} command: lanuch container command
 * @param {subsystems.ResourceConfig} res: resource limitation config
 * @return {*}
 */
func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// create cgroup manager, and call Set and Apply method to limit resource
	cgroupManager := cgroups.NewCgroupManager("Sachet")
	defer cgroupManager.Destroy()

	cgroupManager.Set(res)
	// add container process to cgroup
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(cmdArray, writePipe)
	parent.Wait()

	mntURL := "/home/mnt/"
	rootURL := "/home/"
	container.DeleteWorkSpace(rootURL, mntURL)
	os.Exit(0)
}

/**
 * @description: sent command to child process through pipe
 * @param {[]string} comArray: user command
 * @param {*os.File} writePipe: pipe with child process
 * @return {*}
 */
func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	log.Infof("command all is %s", cmdArray)
	writePipe.WriteString(command)
	writePipe.Close()
}
