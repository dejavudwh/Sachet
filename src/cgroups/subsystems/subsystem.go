/*
 * @Author: dejavudwh
 * @Date: 2021-09-06 23:17:04
 * @LastEditTime: 2021-09-06 23:42:34
 */
package subsystems

/*
	* Structure used to pass resource restrictions
	- MemoryLimit:
	- CpuShare:
	- CpuSet:
*/
type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

/**
 * @description: Subsystems implements these four interfaces to implement resource constraints
 * @method {} Name: return subsystem name
 * @method {path string, pid int} Set: Set the resource limit of a cgroup in this subsystem. The cgroup is the virtual file system path of the hierarchy, so it is abstracted as a path here
 * @method {path string, pid int} Apply: add process to cgroup
 * @method {path string} Remove: remove cgroup
 */
type Subsystem interface {
	Name() string
	Set(path string, res *ResourceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}

// subsystem instance
var SubsystemsIns = []Subsystem{}
