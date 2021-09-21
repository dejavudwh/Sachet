/*
 * @Author: dejavudwh
 * @Date: 2021-09-06 14:08:11
 * @LastEditTime: 2021-09-21 17:12:23
 */
package main

import (
	"Scachet/src/cgroups"
	"Scachet/src/cgroups/subsystems"
	"Scachet/src/container"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

/**
 * @description: Create a new process through the created command
 * @param {bool} tty: IO redirection
 * @param {string} command: lanuch container command
 * @param {subsystems.ResourceConfig} res: resource limitation config
 * @return {*}
 */
func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig, volume string, containerName string) {
	containerID := randStringBytes(10)
	if containerName == "" {
		containerName = containerID
	}

	parent, writePipe := container.NewParentProcess(tty, volume, containerName)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//record container info
	containerName, err := recordContainerInfo(parent.Process.Pid, cmdArray, containerName)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

	// create cgroup manager, and call Set and Apply method to limit resource
	cgroupManager := cgroups.NewCgroupManager("Sachet")
	defer cgroupManager.Destroy()

	cgroupManager.Set(res)
	// add container process to cgroup
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(cmdArray, writePipe)
	if tty {
		parent.Wait()
		deleteContainerInfo(containerName)
	}

	mntURL := "/home/mnt/"
	rootURL := "/home/"
	container.DeleteWorkSpace(rootURL, mntURL, volume)
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

func recordContainerInfo(containerPID int, commandArray []string, containerName string) (string, error) {
	id := randStringBytes(10)
	createTime := time.Now().Format("1900-01-01 11:11:11")
	command := strings.Join(commandArray, "")
	if containerName == "" {
		containerName = id
	}
	containerInfo := &container.ContainerInfo{
		Id:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      container.RUNNING,
		Name:        containerName,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Errorf("Mkdir error %s error %v", dirUrl, err)
		return "", err
	}
	fileName := dirUrl + "/" + container.ConfigName
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %s error %v", fileName, err)
		return "", err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
		return "", err
	}

	return containerName, nil
}

func deleteContainerInfo(containerId string) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerId)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove dir %s error %v", dirURL, err)
	}
}

func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
