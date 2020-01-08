package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var app = cli.NewApp()

var pizza = []string{"Enjoy you test"}

func info() {
	app.Name = "Simple CLI"
	app.Usage = "An example CLI"
	app.Version = "1.0.0"
}

func commands() {
	app.Commands = []*cli.Command{
		{
			Name: "create",
			Aliases: []string{"c"},
			Usage: "Test create",
			Action: func(c *cli.Context) error {
				hello := "hello"
				fmt.Println(hello)
				return nil
			},
		},
		{
			Name: "delete",
			Aliases: []string{"d"},
			Usage: "Test delete",
			Action: func(c *cli.Context) error {
				hello := "world"
				fmt.Println(hello)
				return nil
			},
		},
	}
}

func main() {
	app := cli.NewApp()
	setInformation(app)
	setCommand(app)

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



