package main

import (
	"os"
	"os/signal"

	"github.com/codefresh-io/stevedore/pkg/cmd"
)

const (
	clusterAddedEventName = "cluster:added"
)

func main() {
	handleUnexpectedExit()
	app := cmd.SetupCli()
	app.Run(os.Args)
}

func handleUnexpectedExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			os.Exit(1)
		}
	}()
}
