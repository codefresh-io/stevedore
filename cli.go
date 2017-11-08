package main

import (
	"errors"
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
			Before: func(c *cli.Context) error {
				jwtFromOsEnv := os.Getenv("CODEFRESH_TOKEN")
				if jwtFromOsEnv != "" {
					c.Set("token", jwtFromOsEnv)
				}
				if !c.IsSet("token") {
					fmt.Printf("not set")
					return errors.New("--token nigther CODEFREHS_TOKEN is set")
				}
				return nil
			},
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
