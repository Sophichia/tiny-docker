package main

import (
	"fmt"
	"github.com/Sophichia/tiny-docker/cgroups/subsystems"
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
		var cmdArray []string
		for _, arg := range c.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := c.Bool("ti")
		resCof := &subsystems.ResourceConfig{
			MemoryLimit: c.String("m"),
			CpuShare:    c.String("cpuset"),
			CpuSet:      c.String("cpushare"),
		}

		Run(tty, cmdArray, resCof)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's proocess in container. Do not call it outside",
	Action: func(c *cli.Context) error {
		log.Infof("start init process")
		cmd := c.Args().Get(0)
		log.Info("command %s", cmd)
		err := container.RunContainerInitProcess()
		return err
	},
}
