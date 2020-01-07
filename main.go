package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
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



