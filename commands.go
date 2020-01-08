package main

import (
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
}
