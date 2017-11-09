package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

const (
	success = "SUCCESS"
	failed  = "FAILED"
)

var table = tablewriter.NewWriter(os.Stdout)

func init() {
	table.SetHeader([]string{"Context Name", "Status", "Message"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
}

func addClusterToFinalReport(contextName string, status string, message string) {
	table.Append([]string{contextName, status, message})
}

func reportResult() {
	table.Render()
}
