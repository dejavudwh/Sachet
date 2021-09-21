/*
 * @Author: dejavudwh
 * @Date: 2021-09-21 19:16:53
 * @LastEditTime: 2021-09-21 19:39:26
 */
package main

import (
	"Scachet/src/container"
	_ "Scachet/src/nsenter"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const ENV_EXEC_PID = "sachet_pid"
const ENV_EXEC_CMD = "sachet_cmd"

/**
 * @description: restart container and exec commands
 * @param {string} containerName
 * @param {[]string} comArray
 * @return {*}
 */
func ExecContainer(containerName string, cmdArray []string) {
	pid, err := getContainerPidByName(containerName)
	if err != nil {
		log.Errorf("Exec container getContainerPidByName %s error %v", containerName, err)
		return
	}
	cmdStr := strings.Join(cmdArray, " ")
	log.Infof("container pid %s", pid)
	log.Infof("command %s", cmdStr)

	// re-execute this process to call setns into Namespace
	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// set env variable for callback
	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, cmdStr)

	if err := cmd.Run(); err != nil {
		log.Errorf("Exec container %s error %v", containerName, err)
	}
}

/**
 * @description: get pid by name from container config(log)
 * @param {string} containerName: container name
 * @return {*}
 */
func getContainerPidByName(containerName string) (string, error) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirURL + container.ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return "", err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return "", err
	}

	return containerInfo.Pid, nil
}
