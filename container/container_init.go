package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container gets user command error, cmdArray is nil! ")
	}

	//setUpMount()

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}
	log.Infof("Find path %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Error("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return
	}
	log.Infof("Current location is %s", pwd)
	pivotRoot(pwd)

	// MS_NOEXEC: Not allow other programs in this file system
	// MS_NOSUID: Not allow set-user-ID or set-group-ID when then running process in this file system
	// MS_NODEVV: All mount operation will set it as default after linux2.4
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount root fs to itself error %v", err)
	}
	// Create rootfs/.pivot_root to store old root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	if err := syscall.PivotRoot(root, pivotDIr); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	return os.Remove(pivotDir)
}
