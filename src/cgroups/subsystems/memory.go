/*
 * @Author: dejavudwh
 * @Date: 2021-09-07 00:04:08
 * @LastEditTime: 2021-09-07 00:51:33
 */
package subsystems

import (
	"fmt"
	"io/ioutil"
	"path"
)

type MemorySubSystem struct {
}

/**
 * @description: Set the memory resource limit of the cgroup
 * @param {string} cgroupPath: cgroup (virtual file system)
 * @param {*ResourceConfig} res: restriction information
 * @return {*}
 */
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
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

func (s *MemorySubSystem) Name() string {
	return "memory"
}

// TODO: get subsystem path
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {

	return "", nil
}
