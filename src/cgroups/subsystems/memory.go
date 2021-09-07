/*
 * @Author: dejavudwh
 * @Date: 2021-09-07 00:04:08
 * @LastEditTime: 2021-09-07 11:28:49
 */
package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct {
}

/**
 * @description: Set the memory resource limit of the cgroup
 * @param {string} cgroupPath: cgroup path (virtual file system)
 * @param {*ResourceConfig} res: restriction information
 * @return {*}
 */
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			// After getting the cgroup path, write to the memory limit
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"),
				[]byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			}
		}

		return nil
	} else {
		return err
	}
}

/**
 * @description: Delete the cgroup corresponding to cgroupPath
 * @param {string} cgroupPath: cgroup path (virtual file system)
 * @return {*}
 */
func (s *MemorySubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.Remove(subsysCgroupPath)
	} else {
		return err
	}
}

/**
 * @description: add process to cgroup corresponding to cgroupPath
 * @param {string} cgroupPath: cgroup path (virtual file system)
 * @param {int} pid: process id
 * @return {*}
 */
func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"),
			[]byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}

/**
 * @description: Get subsystem name
 * @param {*}
 * @return {*}
 */
func (s *MemorySubSystem) Name() string {
	return "memory"
}
