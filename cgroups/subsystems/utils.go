package subsystems

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func FindCgroupMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Scan takes text in line by line by default
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		// e.g 37 25 0:30 / /sys/fs/cgroup/memory rw,relatime - cgroup cgroup rw,memory
		fields := strings.Split(text, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ""
}

func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountPoint(subsystem)
	cgroupFullPath := path.Join(cgroupRoot, cgroupPath)
	if _, err := os.Stat(cgroupFullPath); err == nil || (autoCreate && os.IsExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(cgroupFullPath, 0755); err != nil {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return cgroupFullPath, nil
	} else {
		return "", fmt.Errorf("find cgroup path err %v", err)
	}
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}