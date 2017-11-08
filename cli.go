package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func setupCli() *cli.App {
	app := cli.NewApp()
	app.Name = "Stevedore"
	app.Description = "Integrate your connected clusters to your Codefresh account"
	app.Version = "0.0.1"
	setupCommands(app)
	return app
}

var (
	codefreshJwt   string
	kubeConfigPath string
)

func setupCommands(app *cli.App) {
	app.Commands = []cli.Command{
		{
			Name:        "create",
			Description: "Create clusters in Codefresh",
			Action:      create,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "token",
					Usage:       "Codefresh JWT token",
					EnvVar:      "CODEFRESH_TOKEN",
					Destination: &codefreshJwt,
				},
				cli.StringFlag{
					Name:        "config",
					Usage:       "Kubernetes config file to be used as input",
					Value:       fmt.Sprintf("%s/.kube/config", os.Getenv("HOME")),
					EnvVar:      "KUBECONFIG",
					Destination: &kubeConfigPath,
				},
			},
		},
	}
}
