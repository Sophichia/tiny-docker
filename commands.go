package main

import (
	"fmt"
	"github.com/Sophichia/tiny-docker/container"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:      "run",
	Usage:     `Create a container with namespace and cgroup limit.`,
	UsageText: "Example: tiny-docker run -ti [image] [command]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		&cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		&cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		&cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		&cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		&cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
		&cli.StringFlag{
			Name:  "e",
			Usage: "set environment variables",
		},
		&cli.StringFlag{
			Name:  "net",
			Usage: "container network",
		},
		&cli.StringFlag{
			Name:  "p",
			Usage: "port mapping",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Args().Len() < 1 {
			return fmt.Errorf("Missing container command! ")
		}
		cmd := c.Args().Get(0)
		tty := c.Bool("ti")
		Run(tty, cmd)
		return nil
		//
		//var cmdArray []string
		//for _, cmd := range c.Args().Slice() {
		//	cmdArray = append(cmdArray, cmd)
		//}
		//
		//imageName := cmdArray[0]
		//cmdArray = cmdArray[1:]
		//
		//createTty := c.Bool("ti")
		//detach := c.Bool("d")
		//
		//if createTty && detach {
		//	return fmt.Errorf("ti and d param cannot set together! ")
		//}
		//
		//resConfig := &subsystems.ResourceConfig{
		//	MemoryLimit: c.String("m"),
		//	CpuSet: c.String("cpuset"),
		//	CpuShare: c.String("cpushare"),
		//}
		//
		//log.Infof("create tty: %v", createTty)
		//containerName :=  c.String("name")
		//volume := c.String("v")
		//network := c.String("net")
		//
		//envSlice := c.StringSlice("e")
		//portMapping := c.StringSlice("p")
		//
		//err := Run(createTty, cmdArray, resConfig, containerName, volume, imageName, envSlice, network, portMapping)
		//if err != nil {
		//	return fmt.Errorf("run container command fails %v", err)
		//}
		//
		//return nil
	},
}

var initCommand = cli.Command{
	Name: "init",
	Usage: "Init container process run user's proocess in container. Do not call it outside",
	Action: func(c *cli.Context) error {
		log.Infof("start init process")
		cmd := c.Args().Get(0)
		log.Info("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
