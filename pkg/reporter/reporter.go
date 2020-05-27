package reporter

import "fmt"

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
		data []struct {
			name    string
			status  string
			message string
		}
	}
)

func NewReporter() Reporter {
	return &reporter{}
}

func (r *reporter) AddToReport(contextName string, status string, message string) {
	r.data = append(r.data, struct {
		name    string
		status  string
		message string
	}{
		name:    contextName,
		status:  status,
		message: message,
	})
}

func (r *reporter) Print() {
	for _, d := range r.data {
		if d.status == SUCCESS {
			fmt.Printf("Kubernetes context %s added to Codefresh\n", d.name)
			continue
		}

		if d.status == FAILED {
			fmt.Printf("Failed to add Kubernetes context %s to Codefresh.%s\n", d.name, d.message)
			continue
		}
	}
}
