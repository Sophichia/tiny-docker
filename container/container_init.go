package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
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

func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command %s", command)

	// MS_NOEXEC: Not allow other programs in this file system
	// MS_NOSUID: Not allow set-user-ID or set-group-ID when then running process in this file system
	// MS_NODEVV: All mount operation will set it as default after linux2.4
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	argv := []string{command}
	// syscall.Exec will call "int execve" in kernel.
	// It will overwrite all the data and context in stack for init process by the upcoming process.
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf("execute command %v", err.Error())
		return fmt.Errorf("execute command fails %v", err)
	}
	return nil
}
