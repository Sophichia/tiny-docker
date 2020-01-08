package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSubSystem struct {
}

func (s *CpuSubSystem) Set(cgroupPath string, resource *ResourceConfig) error {
	if subSysCroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if resource.CpuShare != "" {
			if err := ioutil.WriteFile(path.Join(subSysCroupPath, "cpu.shares"), []byte(resource.CpuShare), 0644); err != nil {
				return fmt.Errorf("set cgroup cpu share fail %v", err)
			}
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup path failed %v", err)
	}
}

func (s *CpuSubSystem) Remove(cgroupPath string) error {
	if subSysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subSysCgroupPath)
	} else {
		return fmt.Errorf("get cgroup path failed %v", err)
	}
}

func (s *CpuSubSystem) Apply(cgroupPath string, pid int) error {
	if subSysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err = ioutil.WriteFile(path.Join(subSysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("get cgroup path failed %v", err)
	}
}

func (s *CpuSubSystem) Name() string {
	return "cpu"
}
