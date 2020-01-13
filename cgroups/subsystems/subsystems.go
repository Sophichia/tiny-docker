package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type Subsystem interface {
	Name() string // returns the name of subsystem
	Set(path string, res *ResourceConfig) error // set the resource limitation
	Apply(path string, pid int) error // Add process into cgroup
	Remove(path string) error // Delete cgroup
}

var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
