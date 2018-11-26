package reporter

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

const (
	SUCCESS = "SUCCESS"
	FAILED  = "FAILED"
)

type (
	Reporter interface {
		AddToReport(string, string, string)
		Print()
	}

	reporter struct {
		table *tablewriter.Table
	}
)

func NewReporter() Reporter {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Context Name", "Status", "Message"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	return &reporter{
		table: table,
	}
}

func (r *reporter) AddToReport(contextName string, status string, message string) {
	r.table.Append([]string{contextName, status, message})
}

func (r *reporter) Print() {
	r.table.Render()
}
