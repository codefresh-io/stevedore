package stevedore

import (
	"github.com/codefresh-io/stevedore/pkg/codefresh"
	"github.com/codefresh-io/stevedore/pkg/kubernetes"
	"github.com/codefresh-io/stevedore/pkg/reporter"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func Init(c *cli.Context) {
	codefreshAPI := codefresh.NewCodefreshAPI(c.String("api-host"), c.String("token"))
	reporter := reporter.NewReporter()
	kubernetesAPI := kubernetes.NewKubernetesAPI(c.String("config"), codefreshAPI, reporter)
	runOnAllContexts := c.IsSet("all")
	runOnContext := c.String("context")
	if runOnAllContexts {
		kubernetesAPI.GoOverAllContexts()
	} else if runOnContext != "" {
		kubernetesAPI.GoOverContextByName(runOnContext)
	} else {
		kubernetesAPI.GoOverCurrentContext()
	}
	reporter.Print()
	log.Info("Operation is done, check your account setting")
}
