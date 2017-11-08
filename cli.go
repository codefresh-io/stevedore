package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
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
				log.SetLevel(log.WarnLevel)
				log.SetFormatter(&log.TextFormatter{})
				if c.IsSet("verbose") {
					log.SetLevel(log.InfoLevel)
				}

				url := os.Getenv("CODEFRESH_URL")
				if url != "" {
					log.SetLevel(log.DebugLevel)
					log.Debug(fmt.Sprintf("Using other url %s\n", url))
					baseCodefreshURL = url
				}

				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Turn on verbose mode, default is Warning",
				},
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
