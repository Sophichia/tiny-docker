package main

import (
	"github.com/Sophichia/tiny-docker/cgroups"
	"github.com/Sophichia/tiny-docker/cgroups/subsystems"
	"github.com/Sophichia/tiny-docker/container"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// use mydocker-cgroup as cgroup name
	cgroupManager := cgroups.New("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(cmdArray, writePipe)
	parent.Wait()
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

// TODO
// Change it to uuid
//func randStringBytes(n int) string {
//	letterBytes := "1234567890"
//	rand.Seed(time.Now().UnixNano())
//	b := make([]byte, n)
//	for i := range b {
//		b[i] = letterBytes[rand.Intn(len(letterBytes))]
//	}
//	return string(b)
//}
