package main

import (
	"os"
	"os/signal"
)

const (
	clusterAddedEventName = "cluster:added"
)

func main() {
	handleUnexpectedExit()
	app := setupCli()
	app.Run(os.Args)
}

func handleUnexpectedExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			reportResult()
			os.Exit(1)
		}
	}()
}
