package cgroups

import (
	"fmt"
	"github.com/Sophichia/tiny-docker/cgroups/subsystems"
	"github.com/sirupsen/logrus"
)

type CgroupMnanager struct {
	Path     string
	Resource *subsystems.ResourceConfig
}

func New(path string) *CgroupMnanager {
	return &CgroupMnanager{
		Path: path,
	}
}

// Put a process's PID into this cgroup
func (c *CgroupMnanager) Apply(pid int) error {
	for _, subsysIns := range subsystems.SubsystemsIns {
		if err := subsysIns.Apply(c.Path, pid); err != nil {
			return fmt.Errorf("add pid into cgroup fails %v", err)
		}
	}
	return nil
}

// Set the cgroup's resource limitation
func (c *CgroupMnanager) Set(res *subsystems.ResourceConfig) error {
	for _, subsysIns := range subsystems.SubsystemsIns {
		if err := subsysIns.Set(c.Path, res); err != nil {
			return fmt.Errorf("set resource limitation fails %v", err)
		}
	}
	return nil
}

// Destroy cgroup
func (c *CgroupMnanager) Destroy() error {
	for _, subsysIns := range subsystems.SubsystemsIns {
		if err := subsysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fails %v", err)
		}
	}
	return nil
}
