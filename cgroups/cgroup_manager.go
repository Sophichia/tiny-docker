package cgroups

import (
	"fmt"
	"github.com/Sophichia/tiny-docker/cgroups/subsystems"
	"github.com/sirupsen/logrus"
)

type CgroupManager struct {
	Path     string
	Resource *subsystems.ResourceConfig
}

func New(path string) *CgroupManager {
	logrus.Infof("Create cgroupManager %s", path)
	return &CgroupManager{
		Path: path,
	}
}

// Put a process's PID into this cgroup
func (c *CgroupManager) Apply(pid int) error {
	for _, subsysIns := range subsystems.SubsystemsIns {
		if err := subsysIns.Apply(c.Path, pid); err != nil {
			return fmt.Errorf("add pid into cgroup fails %v", err)
		}
	}
	return nil
}

// Set the cgroup's resource limitation
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subsysIns := range subsystems.SubsystemsIns {
		if err := subsysIns.Set(c.Path, res); err != nil {
			return fmt.Errorf("set resource limitation fails %v", err)
		}
	}
	return nil
}

// Destroy cgroup
func (c *CgroupManager) Destroy() error {
	for _, subsysIns := range subsystems.SubsystemsIns {
		if subsystems.FileExists(c.Path) {
			if err := subsysIns.Remove(c.Path); err != nil {
				logrus.Warnf("remove cgroup fails %v", err)
			}
		}
	}
	return nil
}
