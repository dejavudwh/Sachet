/*
 * @Author: dejavudwh
 * @Date: 2021-09-07 10:51:10
 * @LastEditTime: 2021-09-07 11:46:45
 */
package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

/**
 * @description: Find out the root node directory of the hierarchy cgroup where a certain subsystem is mounted through proc/self/mountinfo
 * @param {string} subsystem: subsystem name
 * @return {*}
 */
func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		fields := strings.Split(text, " ")
		// get options
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				// path
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

/**
 * @description: Get the absolute path of the cgroup in the virtual file system
 * @param {string} subsystem
 * @param {string} cgroupPath
 * @param {bool} autoCreate
 * @return {*}
 */
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	// Determine whether the cgroup exists or needs to be created automatically
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err != nil {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
