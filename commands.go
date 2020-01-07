package main

import "github.com/urfave/cli"

var runCommand = cli.Command{
	Name:                   "run",
	Usage:                  `Create a container with namespace and cgroup limit.`,
	UsageText:              "Example: tiny-docker run -ti [image] [command]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "ti",
			Usage: "enable tty",
		},
		&cli.BoolFlag{
			Name: "d",
			Usage: "detach container",
		},
		&cli.StringFlag{
			Name: "m",
			Usage: "memory limit",
		},
	},
}
