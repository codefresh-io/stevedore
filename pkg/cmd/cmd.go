package cmd

import (
	"fmt"
	"os"

	"github.com/codefresh-io/stevedore/stevedore"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func SetupCli() *cli.App {
	app := cli.NewApp()
	app.Name = "Stevedore"
	app.Description = "Integrate your connected clusters to your Codefresh account"
	app.Email = "olegs@gmail.com"
	app.Version = "1.6.2"
	setupCommands(app)
	return app
}

func setupCommands(app *cli.App) {
	app.Commands = []cli.Command{
		{
			Name:        "create",
			Description: "Create clusters in Codefresh. Default is to add current-context",
			Action:      stevedore.Init,
			Before: func(c *cli.Context) error {
				log.SetLevel(log.FatalLevel)
				log.SetFormatter(&log.TextFormatter{})
				if c.IsSet("verbose") {
					log.SetLevel(log.InfoLevel)
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Turn on verbose mode",
				},
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Add all clusters from config file, default is only current context",
				},
				cli.StringFlag{
					Name:  "context, c",
					Usage: "Add spesific cluster",
				},
				cli.StringFlag{
					Name:   "token",
					Usage:  "Codefresh token",
					EnvVar: "CODEFRESH_TOKEN",
				},
				cli.StringFlag{
					Name:   "config",
					Usage:  "Kubernetes config file to be used as input",
					Value:  fmt.Sprintf("%s/.kube/config", os.Getenv("HOME")),
					EnvVar: "KUBECONFIG",
				},
				cli.StringFlag{
					Name:   "api-host",
					Usage:  "Codefresh API host",
					Value:  "https://g.codefresh.io/",
					EnvVar: "CODEFRESH_URL",
				},
				cli.StringFlag{
					Name:   "namespace",
					Usage:  "Which namespace to use while adding cluster to Codefresh (only with --context)",
					Value:  "default",
					EnvVar: "NAMESPACE",
				},
				cli.StringFlag{
					Name:   "serviceaccount",
					Usage:  "Which service account to use while adding cluster to Codefresh (only with --context)",
					Value:  "default",
					EnvVar: "SERVICE_ACCOUNT",
				},
				cli.BoolFlag{
					Name:  "behind-firewall, b",
					Usage: "Spesify whenever the cluster is behined firewall (only with --context)",
				},
				cli.StringFlag{
					Name:   "name-overwrite",
					Usage:  "Spesify under which name save the cluster in Codefresh, default is the same name as the context (only with --context)",
					EnvVar: "NAME_OVERWRITE",
				},
			},
		},
	}
}
