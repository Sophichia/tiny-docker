package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

// The root hierarchy for memory subsystem
const cgroupMemoryHiearchyMount = "/sys/fs/cgroup/memory"

func main() {
	if os.Args[0] == "/proc/self/exe" {
		// container process
		fmt.Printf("current pid %d", syscall.Getpid())
		fmt.Println()
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	} else {
		// Get the forked process pid in outside namespace
		fmt.Printf("%v", cmd.Process.Pid)

		// Create a cgroup in mounted memory subsystem's Hierarchy
		os.Mkdir(path.Join(cgroupMemoryHiearchyMount, "testmemorylimit"), 0755)
		// Add the container process into this cgroup
		ioutil.WriteFile(path.Join(cgroupMemoryHiearchyMount, "testmemorylimit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		// limit the resource
		ioutil.WriteFile(path.Join(cgroupMemoryHiearchyMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644)
		cmd.Process.Wait()
	}
}
