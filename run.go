package main

import (
	"github.com/Sophichia/tiny-docker/container"
	log "github.com/sirupsen/logrus"
	"os"
)

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Errorf("start process fail %v", err)
	}
	parent.Wait()
	os.Exit(-1)
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
