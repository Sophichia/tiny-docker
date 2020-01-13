package main

import (
	"github.com/Sophichia/tiny-docker/cgroups/subsystems"
	"github.com/Sophichia/tiny-docker/container"
	log "github.com/sirupsen/logrus"
	"os"
)

func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig) {
	parent, writePip := container.NewParentProcess()
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
