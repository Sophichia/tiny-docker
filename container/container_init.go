package container

import (
	"os"
	"os/exec"
	"syscall"
)

//func RunContainerInitProcess() error {
//	cmdArray := readUserCommand()
//	if cmdArray == nil || len(cmdArray) == 0 {
//		return fmt.Errorf("get user command error, command array is empty")
//	}
//
//	setUpMount()
//
//	path, err := exec.LookPath(cmdArray[0])
//	if err != nil {
//		log.Errorf("exec loop path error %v", err)
//		return err
//	}
//	log.Infof("find path %s", path)
//	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
//		log.Errorf(err.Error())
//	}
//	return nil
//}
//
//func readUserCommand() []string {
//	pipe := os.NewFile(uintptr(3), "pipe")
//	defer pipe.Close()
//	msg, err := ioutil.ReadAll(pipe)
//	if err != nil {
//		log.Errorf("init reading pipe error %v", err)
//		return nil
//	}
//	msgStr := string(msg)
//	return strings.Split(msgStr, " ")
//}
//
//
//// Set up the mount point
//func setUpMount() {
//	pwd, err := os.Getwd()
//	if err != nil {
//		log.Errorf("Get current location error %v", err)
//		return
//	}
//	log.Infof("Current location is %s", pwd)
//	pivotRoot(pwd)
//
//	// mount proc
//	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
//	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
//
//	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID | syscall.MS_STRICTATIME, "mode=755")
//}

func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}
