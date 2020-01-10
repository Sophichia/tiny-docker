package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	setInformation(app)
	setCommand(app)

	app.Commands = []*cli.Command{
		&initCommand,
		&runCommand,
	}

	app.Before = func(context *cli.Context) error {
		// Log as JSON format instead of ASCII formatter
		log.SetFormatter(&log.JSONFormatter{})

		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func setInformation(app *cli.App) {
	app.Name = "Tiny-docker"
	app.Usage = "Tiny-docker is a simple and naive implementation of container runtime."
	app.Version = "1.0.0"
}

func setCommand(app *cli.App) {
	app.Commands = []*cli.Command{
		&runCommand,
	}
}
